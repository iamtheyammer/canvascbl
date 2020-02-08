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
)

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
