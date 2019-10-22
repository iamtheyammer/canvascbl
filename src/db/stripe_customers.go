package db

import (
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/stripe_customers"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/stripe/stripe-go"
)

func UpsertStripeCustomer(cust stripe.Customer) error {

	return stripe_customers.Upsert(util.DB, &stripe_customers.UpsertRequest{
		StripeID: cust.ID,
		Email:    cust.Email,
	})
}
