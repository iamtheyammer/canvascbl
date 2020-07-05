package gradesapi

import (
	"errors"
	"fmt"
	"github.com/iamtheyammer/canvascbl/backend/src/oauth2"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func OutcomeHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	oID := ps.ByName("outcomeID")
	if len(oID) < 1 || !util.ValidateIntegerString(oID) {
		util.SendBadRequest(w, "missing or invalid outcomeID as url param")
		return
	}

	userID, rdP, sess, errCtx := authorizer(w, r, []oauth2.Scope{oauth2.ScopeOutcomes}, &oauth2.AuthorizerAPICall{
		Method:    "GET",
		RoutePath: "outcomes/:outcomeID",
	})
	if (userID == nil || rdP == nil || errCtx == nil) && sess == nil {
		return
	}

	errCtx.AddCustomField("outcome_id", oID)

	var o *canvasOutcomeResponse
	_, err := handleRequestWithTokenRefresh(func(reqD *requestDetails) error {
		out, outErr := getCanvasOutcome(*reqD, oID)
		if outErr != nil {
			return fmt.Errorf("error getting canvas outcome %s: %w", oID, outErr)
		}

		o = out
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
		handleISE(w, errCtx.Apply(fmt.Errorf("error getting outcome %s: %w", oID, err)))
		return
	}

	go saveOutcomeToDB((*canvasOutcome)(o))

	sendJSON(w, o)
	return
}
