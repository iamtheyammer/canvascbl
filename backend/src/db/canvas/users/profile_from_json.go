package users

import (
	"encoding/json"
	"github.com/pkg/errors"
)

type CanvasProfileResponse struct {
	ID           int64  `json:"id"`
	LoginID      string `json:"login_id"`
	LtiUserID    string `json:"lti_user_id"`
	Name         string `json:"name"`
	PrimaryEmail string `json:"primary_email"`
}

func ProfileFromJSON(pr *string) (*CanvasProfileResponse, error) {
	var p CanvasProfileResponse

	err := json.Unmarshal([]byte(*pr), &p)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal profile JSON")
	}

	return &p, nil
}
