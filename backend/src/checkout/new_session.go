package checkout

import (
	"encoding/json"
	"github.com/iamtheyammer/canvascbl/backend/src/db"
	productssvc "github.com/iamtheyammer/canvascbl/backend/src/db/services/products"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/subscriptions"
	"github.com/iamtheyammer/canvascbl/backend/src/env"
	"github.com/iamtheyammer/canvascbl/backend/src/middlewares"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/checkout/session"
	"net/http"
	"strconv"
	"time"
)

type CreateCheckoutSessionResponse struct {
	Session      string `json:"session"`
	ForProductID uint64 `json:"forProductId"`
}

func CreateCheckoutSessionHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	browserSession := middlewares.Session(w, req, true)
	if browserSession == nil {
		return
	}

	if browserSession.HasValidSubscription {
		util.SendBadRequest(w, "you already have a valid subscription")
		return
	}

	productID := req.URL.Query().Get("productId")
	if len(productID) < 1 {
		util.SendBadRequest(w, "missing productId in query")
		return
	}

	if !util.ValidateIntegerString(productID) {
		util.SendBadRequest(w, "the productId should be an integer")
		return
	}

	pID, err := strconv.Atoi(productID)
	if err != nil {
		util.SendInternalServerError(w)
		return
	}

	cust, err := db.GetStripeCustomer(browserSession.UserID)

	product, err := db.CheckoutListProduct(&productssvc.ListRequest{ID: uint64(pID)})
	if err != nil {
		util.SendInternalServerError(w)
		return
	} else if product == nil {
		util.SendNotFoundWithReason(w, "product not found")
		return
	}

	subs, err := db.GetSubscriptions(&subscriptions.GetRequest{UserID: browserSession.UserID})
	if err != nil {
		util.HandleError(errors.Wrap(err, "error creating checkout session: error getting user subscriptions"))
		util.SendInternalServerError(w)
	}

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		SubscriptionData: &stripe.CheckoutSessionSubscriptionDataParams{
			Items: []*stripe.CheckoutSessionSubscriptionDataItemsParams{
				{
					Plan: stripe.String(product.StripeID),
				},
			},
		},
		SuccessURL: stripe.String(env.StripePurchaseSuccessURL),
		CancelURL:  stripe.String(env.StripeCancelPurchaseURL),
	}

	if cust != nil {
		params.Customer = stripe.String(cust.StripeID)
	} else {
		params.CustomerEmail = stripe.String(browserSession.Email)
	}

	if len(*subs) == 0 {
		// the extra minute is for Stripe-- it seems to use <= so 168 hours is 6 days
		params.SubscriptionData.TrialEnd = stripe.Int64(time.Now().Add(time.Hour * 168).Add(time.Minute).Unix())
	}

	sess, err := session.New(params)
	if err != nil {
		util.HandleError(errors.Wrap(err, "error stripe"))
		util.SendInternalServerError(w)
		return
	}

	resp := CreateCheckoutSessionResponse{
		Session:      sess.ID,
		ForProductID: uint64(pID),
	}

	respJSON, err := json.Marshal(resp)
	if err != nil {
		util.SendInternalServerError(w)
		return
	}

	util.SendJSONResponse(w, respJSON)
}
