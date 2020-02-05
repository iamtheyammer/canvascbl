package oauth2

import (
	"encoding/json"
	"fmt"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/oauth2"
	"github.com/iamtheyammer/canvascbl/backend/src/middlewares"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/url"
	"time"
)

type consentHandlerResponse struct {
	RedirectTo string `json:"redirect_to"`
}

func ConsentHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	consentCode := r.URL.Query().Get("consent_code")
	if len(consentCode) < 1 {
		util.SendBadRequest(w, "missing consent_code as query param")
		return
	} else if !util.ValidateUUIDString(consentCode) {
		util.SendBadRequest(w, "invalid consent_code")
		return
	}

	action := r.URL.Query().Get("action")
	if len(action) < 1 {
		util.SendBadRequest(w, "missing action as query param")
		return
	} else if action != "authorize" && action != "deny" {
		util.SendBadRequest(w, "invalid action as query param")
		return
	}

	sess := middlewares.Session(w, r, true)
	if sess == nil {
		return
	}

	code, err := oauth2.GetCode(util.DB, &oauth2.ListCodesRequest{ConsentCode: consentCode})
	if err != nil {
		util.HandleError(fmt.Errorf("error getting code from consentcode: %w", err))
		util.SendInternalServerError(w)
		return
	}

	if code == nil {
		util.SendBadRequest(w, "invalid consent_code")
		return
	}

	if code.ExpiresAt.Before(time.Now()) {
		util.SendUnauthorized(w, "expired code, restart oauth2 flow")
		return
	}

	redirectURI, err := oauth2.GetRedirectURI(util.DB, &oauth2.ListRedirectURIsRequest{ID: code.RedirectURIID})
	if err != nil {
		util.HandleError(fmt.Errorf("erorr getting redirect uri: %w", err))
		util.SendInternalServerError(w)
		return
	}

	q := url.Values{}

	if action == "authorize" {
		q.Add("code", code.Code)
	} else {
		q.Add("error", "access_denied")
	}

	rURL := redirectURI.RedirectURI + "?" + q.Encode()

	resp := consentHandlerResponse{RedirectTo: rURL}
	jResp, err := json.Marshal(&resp)
	if err != nil {
		util.HandleError(fmt.Errorf("error marshaling consentHandlerResponse: %w", err))
		util.SendInternalServerError(w)
		return
	}

	util.SendJSONResponse(w, jResp)
	return
}
