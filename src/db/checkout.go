package db

import (
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/checkout"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/users"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
	"github.com/stripe/stripe-go"
)

func CheckoutListProducts() (*[]checkout.Product, error) {
	products, err := checkout.ListProducts(util.DB)
	if err != nil {
		return nil, errors.Wrap(err, "error listing products")
	}

	return products, nil
}

func CheckoutListProduct(req *checkout.ListProductRequest) (*checkout.Product, error) {
	return checkout.ListProduct(util.DB, req)
}

func CheckoutWebhookInsertSubscription(req stripe.Subscription) error {
	trx, err := util.DB.Begin()
	if err != nil {
		return errors.Wrap(err, "error beginning insert subscription transaction")
	}

	us, err := users.List(trx, &users.ListRequest{Email: req.Customer.Email})
	if err != nil {
		return errors.Wrap(err, "error listing users")
	}

	if len(*us) < 1 {
		return errors.New("couldn't find a user for the supplied email")
	}

	user := (*us)[0]

	err = checkout.InsertSubscription(util.DB, &checkout.InsertSubscriptionRequest{
		StripeID:           req.ID,
		UserID:             user.ID,
		CustomerStripeID:   req.Customer.ID,
		Plan:               req.Plan.ID,
		Status:             string(req.Status),
		Price:              uint64(req.Plan.Amount),
		CurrentPeriodStart: uint64(req.CurrentPeriodStart),
		CurrentPeriodEnd:   uint64(req.CurrentPeriodEnd),
		TrialEnd:           uint64(req.TrialEnd),
		CanceledAt:         uint64(req.CanceledAt),
	})
	if err != nil {
		return errors.Wrap(err, "error inserting subscription")
	}

	return nil
}

func CheckoutWebhookUpdateSubscription(req stripe.Subscription) error {
	return checkout.UpdateSubscription(util.DB, &checkout.UpdateSubscriptionRequest{
		Identifier: struct {
			StripeID string
			ID       uint64
		}{
			StripeID: req.ID,
		},
		Values: struct {
			UserID             uint64
			Status             string
			Price              uint64
			CurrentPeriodStart uint64
			CurrentPeriodEnd   uint64
			TrialEnd           uint64
			CanceledAt         uint64
		}{
			Status:             string(req.Status),
			Price:              uint64(req.Plan.Amount),
			CurrentPeriodStart: uint64(req.CurrentPeriodStart),
			CurrentPeriodEnd:   uint64(req.CurrentPeriodEnd),
			TrialEnd:           uint64(req.TrialEnd),
			CanceledAt:         uint64(req.CanceledAt),
		},
	})
}
