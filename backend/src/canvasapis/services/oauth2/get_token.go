package oauth2

import (
	"github.com/iamtheyammer/canvascbl/backend/src/env"
	"github.com/iamtheyammer/canvascbl/backend/src/req"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"net/http"
)

// GetAccessFromRedirectResponse gets an access token from a redirect response
func GetAccessFromRedirectResponse(
	code string,
) (*http.Response, string, error) {
	oauth2URL := util.GenerateCanvasURL("/login/oauth2/token")

	q := oauth2URL.Query()
	q.Set("grant_type", "authorization_code")
	q.Set("client_id", env.CanvasOAuth2ClientID)
	q.Set("client_secret", env.CanvasOAuth2ClientSecret)
	q.Set("code", code)
	oauth2URL.RawQuery = q.Encode()

	return req.MakePostRequest(oauth2URL.String())
}

// GetAccessFromRefresh gets an access token from a refresh token
func GetAccessFromRefresh(
	rt string,
) (*http.Response, string, error) {
	oauth2URL := util.GenerateCanvasURL("/login/oauth2/token")

	q := oauth2URL.Query()
	q.Set("grant_type", "refresh_token")
	q.Set("client_id", env.CanvasOAuth2ClientID)
	q.Set("client_secret", env.CanvasOAuth2ClientSecret)
	q.Set("refresh_token", rt)
	oauth2URL.RawQuery = q.Encode()

	return req.MakePostRequest(oauth2URL.String())
}
