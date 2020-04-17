package gradesapi

import (
	"fmt"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/sessions"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/users"
	"github.com/iamtheyammer/canvascbl/backend/src/email"
	"github.com/iamtheyammer/canvascbl/backend/src/env"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	canvasOAuth2AuthURI    = getCanvasOAuth2AuthURI()
	canvasOAuth2ReauthURI  = getCanvasOAuth2ReauthURI()
	canvasOAuth2SuccessURL = func() *url.URL {
		s, err := url.Parse(env.CanvasOAuth2SuccessURI)
		if err != nil {
			panic(fmt.Errorf("errors parsing the canvas oauth2 success url: %w", err))
		}
		return s
	}()
)

type canvasState struct {
	Intent string
	State  string
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

func getCanvasOAuth2AuthURI() string {
	redirectURL := url.URL{
		Host:   env.CanvasDomain,
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
	q.Set("redirect_uri", env.BaseURL+"/api/canvas/oauth2/response")
	q.Set("scope", util.GetScopesList())
	q.Set("state", state)
	redirectURL.RawQuery = q.Encode()

	return redirectURL.String()
}

func getCanvasOAuth2ReauthURI() string {
	redirectURL := url.URL{
		Host:   env.CanvasDomain,
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
	q.Set("redirect_uri", env.BaseURL+"/api/canvas/oauth2/response")
	q.Set("scope", "/auth/userinfo")
	q.Set("state", state)
	redirectURL.RawQuery = q.Encode()

	return redirectURL.String()
}

func getCanvasOAuth2SuccessURI(name string, intent string) string {
	q := canvasOAuth2SuccessURL.Query()
	q.Add("type", "canvas")
	q.Add("name", name)
	q.Add("intent", intent)

	return canvasOAuth2SuccessURL.String() + "?" + q.Encode()
}

// CanvasOAuth2RequestHandler handles forwarding the user to the proper URI for OAuth2 with Canvas.
func CanvasOAuth2RequestHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	intent := r.URL.Query().Get("intent")
	switch intent {
	case "auth":
		util.SendRedirect(w, canvasOAuth2AuthURI)
	case "reauth":
		util.SendRedirect(w, canvasOAuth2ReauthURI)
	default:
		util.SendRedirect(w, canvasOAuth2AuthURI)
	}
}

// CanvasOAuth2ResponseHandler handles the token grant from Canvas.
func CanvasOAuth2ResponseHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

	grantResp, err := getTokenFromAuthorizationCode(code)
	if err != nil {
		handleISE(w, fmt.Errorf("error getting token from authorization code: %w", err))
		return
	}

	if state.Intent == "reauth" {
		profiles, err := users.List(db, &users.ListRequest{CanvasUserID: grantResp.User.ID})
		if err != nil {
			handleISE(w, fmt.Errorf("error listing users in canvas oauth2 response handler (reauth): %w", err))
			return
		}

		if len(*profiles) < 1 {
			// this user is trying to reauth without authing first?
			util.SendRedirect(w, canvasOAuth2AuthURI)
			return
		}

		// this user just wants a new session
		ss, err := sessions.Generate(db, &sessions.GenerateRequest{
			CanvasUserID: grantResp.User.ID,
		})
		if err != nil {
			handleISE(w, fmt.Errorf("error generating session in canvas oauth2 response handler (reauth): %w", err))
		}

		util.AddSessionToResponse(w, *ss)

		util.SendRedirect(w, getCanvasOAuth2SuccessURI(grantResp.User.Name, state.Intent))
		return
	}

	rd := requestDetails{
		Token:        grantResp.AccessToken,
		RefreshToken: grantResp.RefreshToken,
	}

	profile, err := getCanvasProfile(rd, "self")
	if err != nil {
		handleISE(w, fmt.Errorf("error getting canvas self profile in canvas oauth2 response handler: %w", err))
		return
	}

	// doing this synchronously so that we can generate a session
	profileResp, err := users.UpsertProfile(db, &users.UpsertRequest{
		Name:         profile.Name,
		Email:        profile.PrimaryEmail,
		LTIUserID:    profile.LtiUserID,
		CanvasUserID: int64(profile.ID),
	}, true)
	if err != nil {
		handleISE(w, fmt.Errorf("error upserting profile in canvas oauth2 response handler: %w", err))
		return
	}

	if profileResp.InsertedAt.Add(time.Second * 30).After(time.Now()) {
		go email.SendWelcome(profile.PrimaryEmail, profile.Name)
	}

	// this one can be done async
	go saveCanvasOAuth2GrantToDB(grantResp)

	ss, err := sessions.Generate(db, &sessions.GenerateRequest{CanvasUserID: profile.ID})
	if err != nil {
		handleISE(w, fmt.Errorf("error generating session in canvas oauth2 response handler: %w", err))
		return
	}

	util.AddSessionToResponse(w, *ss)

	util.SendRedirect(w, getCanvasOAuth2SuccessURI(profile.Name, state.Intent))

	return
}
