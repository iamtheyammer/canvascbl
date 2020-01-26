package canvasapis

import (
	"encoding/json"
	"fmt"
	"github.com/iamtheyammer/canvascbl/backend/src/canvasapis/services/oauth2"
	"github.com/iamtheyammer/canvascbl/backend/src/canvasapis/services/users"
	"github.com/iamtheyammer/canvascbl/backend/src/db"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/canvas_tokens"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/sessions"
	userssvc "github.com/iamtheyammer/canvascbl/backend/src/db/services/users"
	"github.com/iamtheyammer/canvascbl/backend/src/env"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var OAuth2AuthURI = getOAuth2AuthURI()
var OAuth2ReauthURI = getOAuth2ReauthURI()
var OAuth2SuccessURL = func() *url.URL {
	s, err := url.Parse(env.CanvasOAuth2SuccessURI)
	if err != nil {
		panic(errors.Wrap(err, "errors parsing the canvas oauth2 success url"))
	}
	return s
}()

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

type canvasState struct {
	Intent string
	State  string
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

	state := canvasState{
		Intent: "auth",
		State:  uuid.NewV4().String(),
	}.String()

	q := redirectURL.Query()
	q.Set("client_id", env.CanvasOAuth2ClientID)
	q.Set("response_type", "code")
	q.Set("purpose", purpose)
	q.Set("redirect_uri", env.CanvasOAuth2RedirectURI)
	q.Set("scope", util.GetScopesList())
	q.Set("state", state)
	redirectURL.RawQuery = q.Encode()

	return redirectURL.String()
}

func getOAuth2ReauthURI() string {
	redirectURL := url.URL{
		Host:   fmt.Sprintf("%s.instructure.com", env.CanvasOAuth2Subdomain),
		Path:   "/login/oauth2/auth",
		Scheme: "https",
	}

	purpose := "CanvasCBLAuth"
	switch env.Env {
	case env.EnvironmentStaging:
		purpose += "-staging"
	case env.EnvironmentDevelopment:
		purpose += "-development"
	}

	state := canvasState{
		Intent: "reauth",
		State:  uuid.NewV4().String(),
	}.String()

	q := redirectURL.Query()
	q.Set("client_id", env.CanvasOAuth2ClientID)
	q.Set("response_type", "code")
	q.Set("purpose", purpose)
	q.Set("redirect_uri", env.CanvasOAuth2RedirectURI)
	q.Set("scope", "/auth/userinfo")
	q.Set("state", state)
	redirectURL.RawQuery = q.Encode()

	return redirectURL.String()
}

func getSuccessOAuth2URI(name string, intent string) string {
	q := OAuth2SuccessURL.Query()
	q.Add("type", "canvas")
	q.Add("name", name)
	q.Add("intent", intent)

	return OAuth2SuccessURL.String() + "?" + q.Encode()
}

func (st canvasState) String() string {
	var s string
	s += "intent=" + st.Intent + ";"
	s += "state=" + st.State + ";"
	return s
}

func unmarshalCanvasState(st string) canvasState {
	var s canvasState
	vs := strings.Split(st, ";")
	for _, v := range vs {
		kv := strings.Split(v, "=")
		switch kv[0] {
		case "state":
			s.State = kv[1]
		case "intent":
			s.Intent = kv[1]
		}
	}

	return s
}

/*
OAuth2RequestHandler handles the beginning of the OAuth2 flow with Canvas.
Specifically, it redirects the user to Canvas for permission.
*/
func OAuth2RequestHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	intent := r.URL.Query().Get("intent")
	switch intent {
	case "auth":
		util.SendRedirect(w, OAuth2AuthURI)
	case "reauth":
		util.SendRedirect(w, OAuth2ReauthURI)
	default:
		util.SendRedirect(w, OAuth2AuthURI)
	}
}

func OAuth2ResponseHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	code := r.URL.Query().Get("code")
	state := unmarshalCanvasState(r.URL.Query().Get("state"))

	if len(code) < 1 || len(state.Intent) < 1 || len(state.State) < 1 {
		if len(state.Intent) < 1 || len(state.State) < 1 {
			r.URL.Query().Set("proxy_error", "malformed state")
		}

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

	_, body, err := oauth2.GetAccessFromRedirectResponse(code)
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

	if state.Intent == "reauth" {
		// make sure that they're already a user
		u, err := db.ListUsers(&userssvc.ListRequest{CanvasUserID: uint64(tokenResp.User.ID)})
		if err != nil {
			util.HandleError(errors.Wrap(err, "error listing users for reauth request"))
			util.SendInternalServerError(w)
			return
		}

		if len(*u) == 0 {
			// if they're not already a user we'll send them to the full auth URI
			// this should never happen but hey who the hell knows
			util.SendRedirect(w, OAuth2AuthURI)
			return
		}

		// we just need to make them a session and redirect
		ss, err := db.GenerateSession(&sessions.GenerateRequest{
			CanvasUserID: uint64(tokenResp.User.ID),
		})
		if err != nil {
			util.HandleError(errors.Wrap(err, "error generating a session for a reauth"))
			util.SendInternalServerError(w)
			return
		}

		util.AddSessionToResponse(w, *ss)

		util.SendRedirect(w, getSuccessOAuth2URI(tokenResp.User.Name, state.Intent))
		return
	}

	exp := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	_, pBody, err := users.GetSelfProfile(&util.RequestDetails{
		Token:     tokenResp.AccessToken,
		Subdomain: env.CanvasOAuth2Subdomain,
	})
	if err != nil {
		util.HandleError(errors.Wrap(err, "error getting self profile in oauth2 handler"))
		util.SendInternalServerError(w)
		return
	}

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

	ss, err := db.UpsertProfileAndGenerateSession(&pBody)
	if err != nil {
		util.SendInternalServerError(w)
		return
	}

	util.AddSessionToResponse(w, *ss)

	util.SendRedirect(w, getSuccessOAuth2URI(tokenResp.User.Name, state.Intent))
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
