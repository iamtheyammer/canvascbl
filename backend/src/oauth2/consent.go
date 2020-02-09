package oauth2

import (
	"encoding/json"
	"errors"
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

type dbScope struct {
	Name        string `json:"name"`
	ShortName   string `json:"short_name"`
	Description string `json:"description,omitempty"`
}

type consentTokenHandlerResponse struct {
	ConsentCode    string    `json:"consent_code"`
	CredentialName string    `json:"credential_name"`
	Scopes         []dbScope `json:"scopes"`
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

	err = oauth2.UpdateCode(util.DB, &oauth2.UpdateCodeRequest{
		Where: oauth2.ListCodesRequest{ID: code.ID},
		Set:   oauth2.InsertOAuth2CodeRequest{UserID: &sess.UserID},
	})
	if err != nil {
		util.HandleError(fmt.Errorf("error updating oauth2 code in consent handler: %w", err))
		util.SendInternalServerError(w)
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

// ConsentInfoHandler gets info for a consent token. Session only.
func ConsentInfoHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	consentCode := r.URL.Query().Get("consent_code")
	if len(consentCode) < 1 {
		util.SendBadRequest(w, "missing consent_code as query param")
		return
	} else if !util.ValidateUUIDString(consentCode) {
		util.SendBadRequest(w, "invalid consent_code as query param")
		return
	}

	sess := middlewares.Session(w, r, true)
	if sess == nil {
		return
	}

	code, err := oauth2.GetCode(util.DB, &oauth2.ListCodesRequest{ConsentCode: consentCode})
	if err != nil {
		util.HandleError(fmt.Errorf("error getting oauth2 code: %w", err))
		util.SendInternalServerError(w)
		return
	}

	if code == nil {
		util.SendBadRequest(w, "invalid consent_code")
		return
	}

	cred, err := oauth2.GetCredential(util.DB, &oauth2.ListCredentialsRequest{
		ID:       code.OAuth2CredentialID,
		IsActive: true,
	})
	if err != nil {
		util.HandleError(fmt.Errorf("erorr getting oauth2 credential: %w", err))
		util.SendInternalServerError(w)
		return
	}

	if cred == nil {
		util.HandleError(errors.New("no credential returned from list credentials in consent token handler"))
		util.SendInternalServerError(w)
		return
	}

	scopes, err := oauth2.ListGrantScopes(util.DB, &oauth2.ListGrantScopesRequest{CodeID: code.ID})
	if err != nil {
		util.HandleError(fmt.Errorf("error listing grant scopes in consent token handler: %w", err))
		util.SendInternalServerError(w)
		return
	}

	var s []dbScope
	for _, sc := range *scopes {
		s = append(s, dbScope{
			Name:        sc.Name,
			ShortName:   sc.ShortName,
			Description: sc.Description,
		})
	}

	j, err := json.Marshal(&consentTokenHandlerResponse{
		ConsentCode:    consentCode,
		CredentialName: cred.Name,
		Scopes:         s,
	})
	if err != nil {
		util.HandleError(fmt.Errorf("error marshaling consent token handler response json: %w", err))
		util.SendInternalServerError(w)
		return
	}

	util.SendJSONResponse(w, j)
	return
}
