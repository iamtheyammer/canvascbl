package oauth2

import (
	"github.com/iamtheyammer/canvascbl/backend/src/req"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"net/http"
)

// Delete deletes an access token (for logging out)
func Delete(
	rd *util.RequestDetails,
) (*http.Response, string, error) {
	deleteTokenURL := util.GenerateCanvasURL("/login/oauth2/token")

	q := deleteTokenURL.Query()
	q.Set("access_token", rd.Token)
	deleteTokenURL.RawQuery = q.Encode()

	return req.MakeDeleteRequest(deleteTokenURL.String())
}
