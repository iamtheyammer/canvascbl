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
