package outcomes

import (
	"fmt"
	"github.com/iamtheyammer/canvascbl/backend/src/canvasapis/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"net/http"
)

// GetByID gets an outcome by ID
func GetByID(rd *util.RequestDetails, id string) (*http.Response, string, error) {
	url := fmt.Sprintf(
		"https://%s.instructure.com/api/v1/outcomes/%s",
		rd.Subdomain,
		id,
	)

	return services.MakeAuthenticatedGetRequest(url, rd.Token)
}
