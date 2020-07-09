package db

import (
	"fmt"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/products"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/subscriptions"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/users"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
	"github.com/stripe/stripe-go"
	"time"
)

func CheckoutListProducts() (*[]products.Product, error) {
	products, err := products.ListProducts(util.DB)
	if err != nil {
		return nil, errors.Wrap(err, "error listing products")
	}

	return products, nil
}

func CheckoutListProduct(req *products.ListRequest) (*products.Product, error) {
	return products.ListProduct(util.DB, req)
}

func CheckoutWebhookInsertSubscription(req stripe.Subscription) error {
	user, err := users.GetByStripeID(util.DB, req.Customer.ID)
	if err != nil {
		return errors.Wrap(err, "error listing users")
	}

	if user == nil {
		return errors.New("no customer with that ID")
	}

	err = subscriptions.Insert(util.DB, &subscriptions.InsertRequest{
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

func CheckoutWebhookUpdateSubscription(sub stripe.Subscription) error {
	user, err := users.GetByStripeID(util.DB, sub.Customer.ID)
	if err != nil || user == nil {
		return fmt.Errorf("error listing user by stripe id: %w", err)
	}

	trx, err := util.DB.Begin()
	if err != nil {
		return fmt.Errorf("error beginning checkout webhook update subscription transaction: %w", err)
	}

	err = subscriptions.Update(trx, &subscriptions.UpdateRequest{
		Identifier: struct {
			StripeID string
			ID       uint64
		}{
			StripeID: sub.ID,
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
			Status:             string(sub.Status),
			Price:              uint64(sub.Plan.Amount),
			CurrentPeriodStart: uint64(sub.CurrentPeriodStart),
			CurrentPeriodEnd:   uint64(sub.CurrentPeriodEnd),
			TrialEnd:           uint64(sub.TrialEnd),
			CanceledAt:         uint64(sub.CanceledAt),
		},
	})
	if err != nil {
		err = fmt.Errorf("error updating subscription: %w", err)
		rollbackErr := trx.Rollback()
		if rollbackErr != nil {
			return fmt.Errorf("error rolling back checkout webhook update subscription trx: %w", rollbackErr)
		}
		return err
	}

	subscriptionEnd := time.Unix(sub.CurrentPeriodEnd, 0)
	subscriptionIsValid := subscriptionEnd.After(time.Now()) && (sub.Status == stripe.SubscriptionStatusActive || sub.Status == stripe.SubscriptionStatusTrialing)

	fmt.Printf("subscription: %+v\nsubscriptionIsValid: %t\n", sub, subscriptionIsValid)

	err = users.UpdateUser(trx, &users.UpdateUserRequest{
		WhereID:              user.ID,
		HasValidSubscription: &subscriptionIsValid,
	})
	if err != nil {
		err = fmt.Errorf("error updating user: %w", err)
		rollbackErr := trx.Rollback()
		if rollbackErr != nil {
			return fmt.Errorf("error rolling back checkout webhook update subscription trx at update user: %w", rollbackErr)
		}
		return err
	}

	err = trx.Commit()
	if err != nil {
		err = fmt.Errorf("error committing checkout webhook update subscription trx: %w", err)
		rollbackErr := trx.Rollback()
		if rollbackErr != nil {
			return fmt.Errorf("error rolling back checkout webhook update subscription trx at commit: %w", rollbackErr)
		}
		return err
	}

	return nil
}
