package gradesapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/iamtheyammer/canvascbl/backend/src/oauth2"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func AssignmentsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cID := ps.ByName("courseID")
	if len(cID) < 1 || !util.ValidateIntegerString(cID) {
		util.SendBadRequest(w, "missing or invalid courseID ass url param")
		return
	}

	userID, rdP, sess := authorizer(w, r, []oauth2.Scope{oauth2.ScopeAlignments}, &oauth2.AuthorizerAPICall{
		Method:    "GET",
		RoutePath: "courses/:courseID/assignments",
	})
	if (userID == nil || rdP == nil) && sess == nil {
		return
	}

	var ass *canvasAssignmentsResponse
	_, err := handleRequestWithTokenRefresh(func(reqD *requestDetails) error {
		tAss, outErr := getCanvasCourseAssignments(*reqD, cID)
		if outErr != nil {
			return fmt.Errorf("error getting assignments for course %s: %w", cID, outErr)
		}

		ass = tAss
		return nil
	}, rdP, *userID)
	if errors.Is(err, canvasErrorInvalidAccessTokenError) {
		handleError(w, gradesErrorResponse{
			Error:  gradesErrorRevokedToken,
			Action: gradesErrorActionRedirectToOAuth,
		}, http.StatusForbidden)
		return
	} else if errors.Is(err, canvasErrorInvalidAccessTokenError) {
		handleError(w, gradesErrorResponse{
			Error:  gradesErrorRefreshedTokenError,
			Action: gradesErrorActionRedirectToOAuth,
		}, http.StatusForbidden)
		return
	} else if errors.Is(err, canvasErrorUnknownError) {
		handleError(w, gradesErrorUnknownCanvasErrorResponse, util.CanvasProxyErrorCode)
		return
	} else if err != nil {
		handleISE(w, fmt.Errorf("error getting assignments for course %s: %w", cID, err))
		return
	}

	go saveAssignmentsToDB(*ass, cID)

	jAss, err := json.Marshal(&ass)
	if err != nil {
		handleISE(w, fmt.Errorf("error marshaling assignments for course ID %s: %w", cID, err))
	}

	util.SendJSONResponse(w, jAss)
	return
}
