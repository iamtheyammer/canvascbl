package users

import (
	"fmt"
	"github.com/iamtheyammer/canvascbl/backend/src/req"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"net/http"
)

// GetSelfProfile gets the user's own profile
func GetSelfProfile(rd *util.RequestDetails) (*http.Response, string, error) {
	url := fmt.Sprintf(
		"https://%s.instructure.com/api/v1/users/self/profile",
		rd.Subdomain,
	)
	return req.MakeAuthenticatedGetRequest(url, rd.Token)
}

func GetUserObservees(rd *util.RequestDetails, userID string) (*http.Response, string, error) {
	url := fmt.Sprintf("https://%s.instructure.com/api/v1/users/%s/observees", rd.Subdomain, userID)
	return req.MakeAuthenticatedGetRequest(url, rd.Token)
}
