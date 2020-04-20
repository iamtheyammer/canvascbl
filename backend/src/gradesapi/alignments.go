package gradesapi

import (
	"errors"
	"fmt"
	"github.com/iamtheyammer/canvascbl/backend/src/oauth2"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
)

func AlignmentsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cID := ps.ByName("courseID")
	if len(cID) < 1 || !util.ValidateIntegerString(cID) {
		util.SendBadRequest(w, "missing or invalid courseID as url param")
		return
	}

	sID := r.URL.Query().Get("student_id")
	if len(sID) < 1 {
		util.SendBadRequest(w, "missing or invalid student_id as query param")
		return
	}

	userID, rdP, sess := authorizer(w, r, []oauth2.Scope{oauth2.ScopeAlignments}, &oauth2.AuthorizerAPICall{
		Method:    "GET",
		RoutePath: "courses/:courseID/outcome_alignments",
	})
	if (userID == nil || rdP == nil) && sess == nil {
		return
	}

	var alignments *http.Response
	_, err := handleRequestWithTokenRefresh(func(reqD *requestDetails) error {
		resp, alErr := proxyCanvasOutcomeAlignments(*reqD, cID, sID)
		if alErr != nil {
			return fmt.Errorf("error getting outcome alignments for course %s: %w", cID, alErr)
		}

		alignments = resp
		return nil
	}, rdP, *userID)
	if errors.Is(err, canvasErrorInvalidAccessTokenError) || errors.Is(err, canvasErrorInsufficientScopesOnAccessTokenError) {
		handleError(w, GradesErrorResponse{
			Error:  gradesErrorRevokedToken,
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

	w.Header().Set("Content-Type", "application/json")
	_, err = io.Copy(w, alignments.Body)
	if err != nil {
		handleISE(w, fmt.Errorf("error copying body for outcome alignments: %w", err))
		return
	}

	defer alignments.Body.Close()

	return
}
