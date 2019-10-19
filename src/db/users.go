package db

import (
	"github.com/iamtheyammer/canvascbl/backend/src/db/canvas/users"
	userssvc "github.com/iamtheyammer/canvascbl/backend/src/db/services/users"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
)

// UpsertProfile takes a profile response from Canvas and upserts it into the users table
func UpsertProfile(p *string) {
	cpr, err := users.ProfileFromJSON(p)
	if err != nil {
		handleError(err)
		return
	}

	err = userssvc.UpsertProfile(util.DB, &userssvc.UpsertRequest{
		Name:         cpr.Name,
		Email:        cpr.PrimaryEmail,
		LTIUserID:    cpr.LtiUserID,
		CanvasUserID: cpr.ID,
	})
	if err != nil {
		handleError(errors.Wrap(err, "error upserting profile"))
		return
	}

	return
}
