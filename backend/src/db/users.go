package db

import (
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/subscriptions"
	userssvc "github.com/iamtheyammer/canvascbl/backend/src/db/services/users"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
)

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

func ListObservees(req *userssvc.ListObserveesRequest) (*[]userssvc.Observee, error) {
	return userssvc.ListObservees(util.DB, req)
}
