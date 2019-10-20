package db

import (
	"github.com/iamtheyammer/canvascbl/backend/src/db/canvas/users"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/sessions"
	userssvc "github.com/iamtheyammer/canvascbl/backend/src/db/services/users"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
)

// UpsertProfileAndGenerateSession upserts the user's profile and generates them a session
func UpsertProfileAndGenerateSession(pj *string) (*string, error) {
	cpr, err := users.ProfileFromJSON(pj)
	if err != nil {
		handleError(errors.Wrap(err, "error unmarshaling CanvasProfileResponse"))
		return nil, errors.New("error reading canvas response")
	}

	trx, err := util.DB.Begin()
	if err != nil {
		handleError(errors.Wrap(err, "error beginning upsert profile and generate session trx"))
		return nil, errors.New("error beginning transaction")
	}

	err = userssvc.UpsertProfile(trx, &userssvc.UpsertRequest{
		Name:         cpr.Name,
		Email:        cpr.PrimaryEmail,
		LTIUserID:    cpr.LtiUserID,
		CanvasUserID: cpr.ID,
	})
	if err != nil {
		rollbackErr := trx.Rollback()
		if rollbackErr != nil {
			handleError(errors.Wrap(rollbackErr, "error rolling back upsert profile and "+
				"generate session transaction at upsert profile"))
		}
		handleError(errors.Wrap(err, "error upserting profile"))
		return nil, errors.New("error saving profile")
	}

	ss, err := sessions.Generate(trx, uint64(cpr.ID))
	if err != nil {
		rollbackErr := trx.Rollback()
		if rollbackErr != nil {
			handleError(errors.Wrap(rollbackErr, "error rolling back upsert profile and"+
				"generate session transaction at generate session"))
		}
		handleError(errors.Wrap(err, "error generating session"))
		return nil, errors.New("error generating session")
	}

	err = trx.Commit()
	if err != nil {
		handleError(errors.Wrap(err, "error committing upsert profile and generate session transaction"))
		return nil, errors.New("error saving to database")
	}

	return ss, nil
}
