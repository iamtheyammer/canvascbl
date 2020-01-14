package oauth2

import (
	"github.com/iamtheyammer/canvascbl/backend/src/req"
	"net/http"
)

func GetUserInfo(token string) (*http.Response, string, error) {
	return req.MakeAuthenticatedGetRequest("https://www.googleapis.com/oauth2/v1/userinfo", token)
}
