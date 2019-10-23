package checkout

import (
	"encoding/json"
	"github.com/iamtheyammer/canvascbl/backend/src/db"
	productssvc "github.com/iamtheyammer/canvascbl/backend/src/db/services/products"
	"github.com/iamtheyammer/canvascbl/backend/src/middlewares"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/checkout/session"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type CreateCheckoutSessionResponse struct {
	Session      string `json:"session"`
	ForProductID uint64 `json:"forProductId"`
}

func CreateCheckoutSessionHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	browserSession := middlewares.Session(w, req)
	if browserSession == nil {
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

	email := req.URL.Query().Get("email")
	if len(email) < 1 {
		util.SendBadRequest(w, "missing email in query")
		return
	} else if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		util.SendBadRequest(w, "what was supplied in the email query param doesn't look like an email")
		return
	}

	product, err := db.CheckoutListProduct(&productssvc.ListRequest{ID: uint64(pID)})
	if err != nil {
		util.SendInternalServerError(w)
		return
	} else if product == nil {
		util.SendNotFoundWithReason(w, "product not found")
		return
	}

	trialEnd := time.Now().Add(time.Hour * 168).Unix()

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		CustomerEmail: stripe.String(email),
		SubscriptionData: &stripe.CheckoutSessionSubscriptionDataParams{
			Items: []*stripe.CheckoutSessionSubscriptionDataItemsParams{
				{
					Plan: stripe.String(product.StripeID),
				},
			},
			TrialEnd: &trialEnd,
		},
		SuccessURL: stripe.String("http://localhost:3000/#/dashboard/checkout/thanks"),
		CancelURL:  stripe.String("http://localhost:3000/#/dashboard/checkout"),
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
