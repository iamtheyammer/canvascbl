package oauth2

import (
	"encoding/json"
	"github.com/iamtheyammer/canvascbl/backend/src/req"
	"github.com/pkg/errors"
	"net/http"
)

type googleAuthPayload struct {
	Code         string `json:"code"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
	GrantType    string `json:"grant_type"`
}

func GetAccessFromRedirect(
	code string,
	clientID string,
	clientSecret string,
	redirectURI string,
) (*http.Response, string, error) {
	j := googleAuthPayload{
		Code:         code,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
		GrantType:    "authorization_code",
	}

	jb, err := json.Marshal(j)
	if err != nil {
		return nil, "", errors.Wrap(err, "error marshaling google request token json payload")
	}

	return req.MakePostRequestWithBody("https://oauth2.googleapis.com/token", jb)
}
