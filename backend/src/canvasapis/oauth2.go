package canvasapis

import (
	"encoding/json"
	"fmt"
	"github.com/iamtheyammer/canvascbl/backend/src/canvasapis/services/oauth2"
	"github.com/iamtheyammer/canvascbl/backend/src/canvasapis/services/users"
	"github.com/iamtheyammer/canvascbl/backend/src/db"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/canvas_tokens"
	"github.com/iamtheyammer/canvascbl/backend/src/env"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"time"
)

var OAuth2AuthURI = getOAuth2AuthURI()

type canvasTokenGrantResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	User         struct {
		EffectiveLocale string `json:"effective_locale"`
		GlobalID        string `json:"global_id"`
		ID              int64  `json:"id"`
		Name            string `json:"name"`
	} `json:"user"`
}

func getOAuth2AuthURI() string {
	redirectURL := url.URL{
		Host:   fmt.Sprintf("%s.instructure.com", env.CanvasOAuth2Subdomain),
		Path:   "/login/oauth2/auth",
		Scheme: "https",
	}

	purpose := "CanvasCBL"
	switch env.Env {
	case env.EnvironmentStaging:
		purpose += "-staging"
	case env.EnvironmentDevelopment:
		purpose += "-development"
	}

	q := redirectURL.Query()
	q.Set("client_id", env.CanvasOAuth2ClientID)
	q.Set("response_type", "code")
	q.Set("purpose", purpose)
	q.Set("redirect_uri", env.CanvasOAuth2RedirectURI)
	q.Set("scope", util.GetScopesList())
	redirectURL.RawQuery = q.Encode()

	return redirectURL.String()
}

/*
OAuth2RequestHandler handles the beginning of the OAuth2 flow with Canvas.
Specifically, it redirects the user to Canvas for permission.
*/
func OAuth2RequestHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	util.SendRedirect(w, OAuth2AuthURI)
}

func OAuth2ResponseHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	code := r.URL.Query().Get("code")

	if len(code) < 1 {
		// an error occurred, just redirect
		r.URL.Query().Set("error_source", "canvas")
		util.SendRedirect(
			w,
			fmt.Sprintf(
				"%s?%s",
				env.CanvasOAuth2SuccessURI,
				r.URL.Query().Encode(),
			),
		)
		return
	}

	resp, body, err := oauth2.GetAccessFromRedirectResponse(code)
	if err != nil {
		util.SendInternalServerError(w)
		return
	}

	var tokenResp canvasTokenGrantResponse
	err = json.Unmarshal([]byte(body), &tokenResp)
	if err != nil {
		util.HandleError(errors.Wrap(err, "error decoding canvas token grant response"))
		util.SendInternalServerError(w)
		return
	}

	exp := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	err = db.InsertCanvasToken(&canvas_tokens.InsertRequest{
		CanvasUserID: uint64(tokenResp.User.ID),
		Token:        tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		ExpiresAt:    &exp,
	})
	if err != nil {
		util.HandleError(errors.Wrap(err, "error inserting canvas token"))
		util.SendInternalServerError(w)
		return
	}

	_, pBody, err := users.GetSelfProfile(&util.RequestDetails{
		Token:     tokenResp.AccessToken,
		Subdomain: env.CanvasOAuth2Subdomain,
	})
	if err != nil {
		util.HandleError(errors.Wrap(err, "error getting self profile in oauth2 handler"))
		util.SendInternalServerError(w)
		return
	}

	ss, err := db.UpsertProfileAndGenerateSession(&pBody)
	if err != nil {
		util.SendInternalServerError(w)
		return
	}

	util.AddSessionToResponse(w, *ss)

	util.HandleCanvasOAuth2Response(w, resp, body)
	return
}

func OAuth2RefreshTokenHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	refreshToken := r.URL.Query().Get("refresh_token")
	if len(refreshToken) < 1 {
		util.SendBadRequest(w, "no refresh_token specified")
		return
	}

	resp, body, err := oauth2.GetAccessFromRefresh(refreshToken)
	if err != nil {
		util.SendInternalServerError(w)
		return
	}

	util.HandleCanvasResponse(w, resp, body)
	return
}

func DeleteOAuth2TokenHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ok, rd := util.GetRequestDetailsFromRequest(r)

	if !ok {
		util.SendUnauthorized(w, util.RequestDetailsFailedValidationMessage)
		return
	}

	resp, body, err := oauth2.Delete(rd)
	if err != nil {
		util.SendInternalServerError(w)
		return
	}

	util.HandleCanvasResponse(w, resp, body)
	return
}
