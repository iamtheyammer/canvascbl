package db

import (
	"github.com/iamtheyammer/canvascbl/backend/src/db/canvas/users"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/subscriptions"
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

func GetUserFromCanvasProfileResponseJSON(p *string) *userssvc.User {
	cpr, err := users.ProfileFromJSON(p)
	if err != nil {
		handleError(errors.Wrap(err, "error getting users from json"))
		return nil
	}

	usP, err := userssvc.ListFromCanvasResponse(util.DB, cpr)
	if err != nil {
		handleError(errors.Wrap(err, "error getting users from canvas response from json"))
		return nil
	}

	us := *usP

	if len(us) != 1 {
		handleError(errors.New("not 1 user returned from get users from canvas response from json"))
	}

	return &us[0]
}

func GetUserFromStripeSubscriptionID(stripeSubscriptionID string) *userssvc.User {
	subsP, err := subscriptions.Get(util.DB, &subscriptions.GetRequest{StripeID: stripeSubscriptionID})
	if err != nil {
		handleError(errors.Wrap(err, "error getting subscription from stripe subscription id"))
		return nil
	}

	subs := *subsP
	if len(subs) < 1 {
		return nil
	}

	user, err := userssvc.GetByStripeID(util.DB, subs[0].CustomerStripeID)
	if err != nil {
		handleError(errors.Wrap(err, "error getting user by stripe ID"))
		return nil
	}

	return user
}

func ListUsers(req *userssvc.ListRequest) (*[]userssvc.User, error) {
	return userssvc.List(util.DB, req)
}

func HandleObservees(observees *string, requestingUserID uint64) {
	obsP, err := users.ObserveesFromJSON(observees)
	if err != nil {
		handleError(errors.Wrap(err, "error unmarshaling observees"))
		return
	}

	obs := *obsP

	// if the user has no observees, let's do nothing.
	if len(obs) < 1 {
		return
	}

	// now, we'll start a db transaction
	trx, err := util.DB.Begin()
	if err != nil {
		handleError(errors.Wrap(err, "error beginning handle observees transaction"))
		return
	}

	// get the user's current observees
	dbObserveesP, err := userssvc.ListObservees(trx, &userssvc.ListObserveesRequest{ObserverCanvasUserID: requestingUserID})
	if err != nil {
		handleError(errors.Wrap(err, "error listing user observees"))
		return
	}

	dbObservees := *dbObserveesP

	var (
		toSoftDelete, toUnSoftDelete []uint64
		toUpsert                     []userssvc.Observee
	)

	for _, o := range obs {
		foundIDMatch := false
		for _, dbO := range dbObservees {
			if o.ID == dbO.CanvasUserID {
				// if the names don't match, upsert
				if o.Name != dbO.Name {
					toUpsert = append(toUpsert, userssvc.Observee{
						CanvasUserID: o.ID,
						Name:         o.Name,
					})
				}

				// if it was previously deleted, undelete
				if !dbO.DeletedAt.IsZero() {
					toUnSoftDelete = append(toUnSoftDelete, dbO.CanvasUserID)
				}

				foundIDMatch = true
			}
		}

		// if it exists in observees from canvas but not in db, upsert.
		if !foundIDMatch {
			toUpsert = append(toUpsert, userssvc.Observee{
				CanvasUserID: o.ID,
				Name:         o.Name,
			})
		}
	}

	for _, dbO := range dbObservees {
		foundIDMatch := false
		for _, o := range obs {
			if dbO.CanvasUserID == o.ID {
				foundIDMatch = true
			}
		}

		// if it exists in the db and it's not already deleted
		if !foundIDMatch && dbO.DeletedAt.IsZero() {
			toSoftDelete = append(toSoftDelete, dbO.CanvasUserID)
		}
	}

	if len(toSoftDelete) > 0 {
		err := userssvc.SoftDeleteUserObservees(trx, toSoftDelete)
		if err != nil {
			handleError(errors.Wrap(err, "error soft deleting user observees"))
			rollbackErr := trx.Rollback()
			if rollbackErr != nil {
				handleError(errors.Wrap(err, "error rolling back handle observees trx"))
			}
			return
		}
	}

	if len(toUnSoftDelete) > 0 {
		err := userssvc.UnSoftDeleteUserObservees(trx, toUnSoftDelete)
		if err != nil {
			handleError(errors.Wrap(err, "error un-soft deleting user observees"))
			rollbackErr := trx.Rollback()
			if rollbackErr != nil {
				handleError(errors.Wrap(err, "error rolling back handle observees trx"))
			}
			return
		}
	}

	if len(toUpsert) > 0 {
		err := userssvc.UpsertUserObservees(trx, &userssvc.UpsertObserveesRequest{
			Observees:            toUpsert,
			ObserverCanvasUserID: requestingUserID,
		})
		if err != nil {
			handleError(errors.Wrap(err, "error upserting user observees"))
			rollbackErr := trx.Rollback()
			if rollbackErr != nil {
				handleError(errors.Wrap(err, "error rolling back handle observees trx"))
			}
			return
		}
	}

	err = trx.Commit()
	if err != nil {
		handleError(errors.Wrap(err, "error committing handle observees trx"))
		rollbackErr := trx.Rollback()
		if rollbackErr != nil {
			handleError(errors.Wrap(err, "error rolling back handle observees trx"))
		}
		return
	}

	return
}

func ListObservees(req *userssvc.ListObserveesRequest) (*[]userssvc.Observee, error) {
	return userssvc.ListObservees(util.DB, req)
}
