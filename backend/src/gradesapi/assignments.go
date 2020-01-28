package gradesapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/sessions"
	"github.com/iamtheyammer/canvascbl/backend/src/middlewares"
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

	session := middlewares.Session(w, r)
	if session == nil {
		return
	}

	if session.Type == sessions.VerifiedSessionTypeSessionString {
		// sick
	} else if session.Type == sessions.VerifiedSessionTypeAPIKey {
		util.SendUnauthorized(w, "api keys aren't implemented yet")
		return
	} else {
		util.SendUnauthorized(w, "unsupported authentication method")
		return
	}

	rd, err := rdFromCanvasUserID(session.CanvasUserID)
	if err != nil {
		handleISE(w, fmt.Errorf("error getting rd from canvas user id: %w", err))
		return
	}

	if rd.TokenID < 1 {
		handleError(w, gradesErrorResponse{
			Error:  gradesErrorNoTokens,
			Action: gradesErrorActionRedirectToOAuth,
		}, http.StatusForbidden)
		return
	}

	var ass *canvasAssignmentsResponse
	rd, err = handleRequestWithTokenRefresh(func(reqD *requestDetails) error {
		tAss, outErr := getCanvasCourseAssignments(*reqD, cID)
		if outErr != nil {
			return fmt.Errorf("error getting assignments for course %s: %w", cID, outErr)
		}

		ass = tAss
		return nil
	}, &rd, session.CanvasUserID)
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
