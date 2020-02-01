package gradesapi

import (
	"errors"
	"fmt"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/sessions"
	"github.com/iamtheyammer/canvascbl/backend/src/middlewares"
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

	session := middlewares.Session(w, r, true)
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

	var alignments *http.Response
	rd, err = handleRequestWithTokenRefresh(func(reqD *requestDetails) error {
		resp, alErr := proxyCanvasOutcomeAlignments(*reqD, cID, sID)
		if alErr != nil {
			return fmt.Errorf("error getting outcome alignments for course %s: %w", cID, alErr)
		}

		alignments = resp
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

	w.Header().Set("Content-Type", "application/json")
	_, err = io.Copy(w, alignments.Body)
	if err != nil {
		handleISE(w, fmt.Errorf("error copying body for outcome alignments: %w", err))
		return
	}

	defer alignments.Body.Close()

	return
}
