package middlewares

import (
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"net/http"
	"strings"
)

/*
Bearer just returns the OAuth2 Token for a request.
It's in the Authorization header, after "Bearer ".

It returns the access token (if present) along with a bool representing whether,
if a token is supplied, it is valid.
*/
func Bearer(w http.ResponseWriter, r *http.Request, shouldRespond bool) (string, bool) {
	accessToken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	if len(accessToken) < 1 {
		if shouldRespond {
			util.SendUnauthorized(w, "missing access token")
		}

		return "", true
	} else if !util.ValidateUUIDString(accessToken) {
		if shouldRespond {
			util.SendUnauthorized(w, "invalid access token")
		}

		return "", false
	}

	return accessToken, true
}
