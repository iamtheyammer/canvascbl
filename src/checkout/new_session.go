package checkout

import (
	"fmt"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/checkout/session"
	"net/http"
)

func CreateCheckoutSessionHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		SubscriptionData: &stripe.CheckoutSessionSubscriptionDataParams{
			Items: []*stripe.CheckoutSessionSubscriptionDataItemsParams{
				{
					Plan: stripe.String("plan_G1aF4tORdb3dmG"),
				},
			},
		},
		SuccessURL: stripe.String("http://localhost:3000/#/dashboard/checkout/thanks"),
		CancelURL:  stripe.String("http://localhost:3000/#/dashboard/checkout"),
	}

	sess, err := session.New(params)
	if err != nil {
		fmt.Println(errors.Wrap(err, "error stripe"))
		util.SendInternalServerError(w)
		return
	}

	_, _ = fmt.Fprint(w, sess.ID)
}
