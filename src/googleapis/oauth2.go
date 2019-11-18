package googleapis

import (
	"encoding/json"
	"github.com/iamtheyammer/canvascbl/backend/src/db"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/canvas_tokens"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/google_users"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/sessions"
	userssvc "github.com/iamtheyammer/canvascbl/backend/src/db/services/users"
	"github.com/iamtheyammer/canvascbl/backend/src/env"
	"github.com/iamtheyammer/canvascbl/backend/src/googleapis/services/oauth2"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
)

type googleOAuth2TokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

type googleProfileResponse struct {
	Email string `json:"email"`
	// last name
	FamilyName string `json:"family_name"`
	// optional
	Gender string `json:"gender"`
	// first name
	GivenName string `json:"given_name"`
	// only present on G Suite accounts
	HostedDomain string `json:"hd"`
	ID           string `json:"id"`
	// Google+ link; only present on G Suite accounts
	Link   string `json:"link"`
	Locale string `json:"locale"`
	// GivenName FamilyName
	Name string `json:"name"`
	// profile picture URL
	Picture       string `json:"picture"`
	VerifiedEmail bool   `json:"verified_email"`
}

var oAuth2RequestURI = getOAuth2RequestURI()

func getOAuth2RequestURI() string {
	rURI := url.URL{
		Scheme: "https",
		Opaque: "",
		User:   nil,
		Host:   "accounts.google.com",
		Path:   "/o/oauth2/v2/auth",
	}

	q := rURI.Query()
	q.Set("client_id", env.GoogleOAuth2ClientID)
	q.Set("redirect_uri", env.GoogleOAuth2RedirectURI)
	q.Set("scope", util.GetGoogleScopesList())
	q.Set("response_type", "code")
	q.Set("prompt", "select_account")
	rURI.RawQuery = q.Encode()

	return rURI.String()
}

func OAuth2RequestHandler(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	util.SendRedirect(w, oAuth2RequestURI)
}

func OAuth2ResponseHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	redirectTo, err := url.Parse(env.CanvasOAuth2SuccessURI)
	if err != nil {
		util.SendInternalServerError(w)
		return
	}

	q := redirectTo.Query()
	q.Set("type", "google")

	// sr sends response
	sr := func() {
		redirectTo.RawQuery = q.Encode()

		util.SendRedirect(w, redirectTo.String())
		return
	}

	qErr := r.URL.Query().Get("error")
	if len(qErr) > 1 {
		q.Set("error", "proxy_google_error")
		q.Set("error_source", "google")
		q.Set("google_error", qErr)

		sr()
		return
	}

	code := r.URL.Query().Get("code")
	if len(code) < 1 {
		util.SendBadRequest(w, "missing code as query param")
		return
	}

	// get an access token for the user
	resp, body, err := oauth2.GetAccessFromRedirect(
		code,
		env.GoogleOAuth2ClientID,
		env.GoogleOAuth2ClientSecret,
		env.GoogleOAuth2RedirectURI,
	)
	if err != nil || resp.StatusCode < 200 || resp.StatusCode > 299 {
		if err == nil {
			util.HandleError(errors.New("error getting google oauth2 access token from code (redirect): status code not 200-299"))
		} else {
			util.HandleError(errors.Wrap(err, "error getting google oauth2 access token from code (redirect)"))
		}
		q.Set("error", "proxy_google_error")
		q.Set("error_source", "google")
		if err == nil {
			q.Set("body", body)
		}

		sr()
		return
	}

	var tr googleOAuth2TokenResponse
	err = json.Unmarshal([]byte(body), &tr)
	if err != nil {
		util.HandleError(errors.Wrap(err, "error unmarshaling google access token response"))
		util.SendInternalServerError(w)
		return
	}

	// get user identity

	resp, body, err = oauth2.GetUserInfo(tr.AccessToken)
	if err != nil || resp.StatusCode < 200 || resp.StatusCode > 299 {
		if err == nil {
			util.HandleError(errors.New("error getting google oauth2 user info: status code not 200-299"))
		} else {
			util.HandleError(errors.Wrap(err, "error getting google oauth2 user info"))
		}
		q.Set("error", "proxy_google_error")
		q.Set("error_source", "google")
		if err == nil {
			q.Set("body", body)
		}

		sr()
		return
	}

	// user profile
	var up googleProfileResponse
	err = json.Unmarshal([]byte(body), &up)
	if err != nil {
		util.HandleError(errors.Wrap(err, "error unmarshaling google profile response"))
		util.SendInternalServerError(w)
		return
	}

	if !util.ValidateGoogleHD(up.HostedDomain) {
		q.Set("error", "proxy_google_error")
		q.Set("error_source", "canvas_proxy")
		q.Set("error_text", "domain not allowed")

		sr()
		return
	}

	// get user by email (do they already have an account?)

	usersP, err := db.ListUsers(&userssvc.ListRequest{Email: up.Email})
	if err != nil {
		util.HandleError(errors.Wrap(err, "error listing users in google oauth2 response"))
		util.SendInternalServerError(w)
		return
	}
	us := *usersP

	ur := google_users.UpsertRequest{
		GoogleID:          up.ID,
		Email:             up.Email,
		GivenName:         &up.GivenName,
		FamilyName:        &up.FamilyName,
		Name:              &up.Name,
		ProfilePictureURL: &up.Picture,
		Gender:            &up.Gender,
		HostedDomain:      &up.HostedDomain,
	}

	if len(us) == 1 {
		ur.UserID = &us[0].ID
	}

	googleUserID, err := db.UpsertGoogleProfile(&ur)
	if err != nil {
		util.HandleError(errors.Wrap(err, "error upserting google profile"))
		util.SendInternalServerError(w)
		return
	}

	sgr := sessions.GenerateRequest{
		GoogleUsersID: *googleUserID,
	}

	if len(us) > 0 {
		sgr.CanvasUserID = us[0].CanvasUserID
	}

	ss, err := db.GenerateSession(&sgr)
	if err != nil {
		util.HandleError(errors.Wrap(err, "error generating session after google oauth2 response"))
		util.SendInternalServerError(w)
		return
	}

	util.AddSessionToResponse(w, *ss)

	// can be overwritten
	q.Set("has_token", "false")

	if len(us) > 0 {
		ctsP, err := db.ListCanvasTokens(&canvas_tokens.ListRequest{
			UserID:         us[0].ID,
			NotExpiredOnly: true,
			Limit:          1,
		})
		if err != nil {
			util.HandleError(errors.Wrap(err, "error listing canvas tokens after google oauth2 response"))
			util.SendInternalServerError(w)
			return
		}

		cts := *ctsP

		if len(cts) > 0 {
			q.Set("has_token", "true")
		}
	}

	sr()
	return
}
