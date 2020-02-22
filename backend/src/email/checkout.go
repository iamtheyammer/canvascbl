package email

import (
	"fmt"
	"github.com/iamtheyammer/canvascbl/backend/src/db"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/products"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
	"github.com/stripe/stripe-go"
	"strings"
)

func SendPurchaseAcknowledgement(sub *stripe.Subscription) {
	user := db.GetUserFromStripeSubscriptionID(sub.ID)
	if user == nil {
		util.HandleError(errors.New("error sending purchase acknowledgement: error getting user from stripe subscription ID"))
		return
	}

	prod, err := db.CheckoutListProduct(&products.ListRequest{
		StripeID: sub.Plan.ID,
	})
	if err != nil {
		util.HandleError(errors.Wrap(err, "error sending purchase acknowledgement: error listing products by stripe id"))
		return
	}
	if prod == nil {
		util.HandleError(errors.New("error sending purchase acknowledgement: no products returned"))
		return
	}

	send(
		purchaseAcknowledgement,
		map[string]interface{}{
			"first_name":   strings.Split(user.Name, " ")[0],
			"product_name": prod.Name,
			"price":        fmt.Sprintf("$%.2f", float64(sub.Plan.Amount)/float64(100)),
		},
		user.Email,
		user.Name,
	)
}

func SendCancellationAcknowledgement(sub *stripe.Subscription) {
	user := db.GetUserFromStripeSubscriptionID(sub.ID)
	if user == nil {
		util.HandleError(errors.New("error sending cancellation acknowledgement: error listing products by stripe id"))
		return
	}

	send(
		cancellationAcknowledgement,
		map[string]interface{}{
			"first_name": strings.Split(user.Name, " ")[0],
		},
		user.Email,
		user.Name,
	)
}
