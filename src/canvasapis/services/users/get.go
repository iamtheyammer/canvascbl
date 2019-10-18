package users

import (
	"fmt"
	"github.com/iamtheyammer/canvascbl/backend/src/canvasapis/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"net/http"
)

// GetSelfProfile gets the user's own profile
func GetSelfProfile(rd *util.RequestDetails) (*http.Response, string, error) {
	url := fmt.Sprintf(
		"https://%s.instructure.com/api/v1/users/self/profile",
		rd.Subdomain,
	)
	return services.MakeAuthenticatedGetRequest(url, rd.Token)
}
