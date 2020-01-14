package checkout

import (
	"encoding/json"
	"fmt"
	"github.com/iamtheyammer/canvascbl/backend/src/db"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/gift_cards"
	productssvc "github.com/iamtheyammer/canvascbl/backend/src/db/services/products"
	"github.com/iamtheyammer/canvascbl/backend/src/middlewares"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/sub"
	"net/http"
	"time"
)

type redeemHandlerBody struct {
	Codes []string `json:"codes"`
}

type redeemHandlerResponse struct {
	Success               bool   `json:"success"`
	Error                 string `json:"error"`
	SubscriptionExpiresAt int64  `json:"subscription_expires_at"`
}

func RedeemHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	session := middlewares.Session(w, r)
	if session == nil {
		return
	}

	if session.HasValidSubscription {
		util.SendBadRequest(w, "you already have a valid subscription")
		return
	}

	var body redeemHandlerBody
	err := middlewares.DecodeJSONBody(r.Body, &body)
	if err != nil {
		util.SendBadRequest(w, "malformed body")
		return
	}

	if len(body.Codes) < 1 {
		util.SendBadRequest(w, "no codes")
		return
	}

	if len(body.Codes) > 10 {
		util.SendBadRequest(w, "you can only send a max of 10 codes at a time")
		return
	}

	for _, code := range body.Codes {
		if !util.ValidateGiftCardClaimCode(code) {
			util.SendBadRequest(w, fmt.Sprintf("this code does not look like a gift card claim code: %s", code))
			return
		}
	}

	cardsP, err := db.ListGiftCards(&gift_cards.ListRequest{
		ClaimCodes: body.Codes,
		ValidOnly:  true,
	})
	if err != nil {
		util.HandleError(errors.Wrap(err, "error listing gift cards in redeemhandler"))
		util.SendInternalServerError(w)
		return
	}
	cards := *cardsP

	if len(cards) != len(body.Codes) {
		util.SendBadRequest(w, "one or more of your codes is invalid")
		return
	}

	product, err := db.CheckoutListProduct(&productssvc.ListRequest{
		ShortName: "plus_monthly",
	})
	if err != nil {
		util.HandleError(errors.Wrap(err, "error listing products in redeemhandler"))
		util.SendInternalServerError(w)
		return
	}

	if product == nil {
		util.HandleError(errors.Wrap(err, "no product with short_name = plus_monthly in redeemhandler"))
		util.SendInternalServerError(w)
		return
	}

	stripeCust, err := db.GetStripeCustomer(session.UserID)
	if err != nil {
		util.HandleError(errors.Wrap(err, "error listing stripe customer in redeemhandler"))
		util.SendInternalServerError(w)
		return
	}

	var stripeCustID string

	if stripeCust == nil {
		cust, err := customer.New(&stripe.CustomerParams{
			Email: stripe.String(session.Email),
		})
		if err != nil {
			util.HandleError(errors.Wrap(err, "error creating stripe customer in redeemhandler"))
			util.SendInternalServerError(w)
			return
		}

		err = db.UpsertStripeCustomer(*cust)
		if err != nil {
			util.HandleError(errors.Wrap(err, "error upserting stripe customer in redeemhandler"))
			util.SendInternalServerError(w)
			return
		}

		stripeCustID = cust.ID
	} else {
		stripeCustID = stripeCust.StripeID
	}

	subscriptionEndDate := time.Now()
	var cardIDs []uint64

	for _, gc := range cards {
		subscriptionEndDate = subscriptionEndDate.Add(time.Duration(gc.ValidFor) * time.Second)
		cardIDs = append(cardIDs, gc.ID)
	}

	subscriptionEndUnix := subscriptionEndDate.Unix()

	subscriptionParams := &stripe.SubscriptionParams{
		CancelAt: stripe.Int64(subscriptionEndUnix),
		Customer: stripe.String(stripeCustID),
		TrialEnd: stripe.Int64(subscriptionEndUnix),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Plan: stripe.String(product.StripeID),
			},
		},
	}

	_, err = sub.New(subscriptionParams)
	if err != nil {
		util.HandleError(errors.Wrap(err, "error creating stripe subscription in redeemhandler"))
		util.SendInternalServerError(w)
		return
	}

	err = db.UpdateGiftCards(&gift_cards.UpdateRequest{
		Where: gift_cards.ListRequest{
			IDs: cardIDs,
		},
		Set: gift_cards.GiftCard{
			RedeemedBy: session.UserID,
			RedeemedAt: time.Now(),
		},
	})
	if err != nil {
		util.HandleError(errors.Wrap(err, "error updating gift cards in redeemhandler"))
		util.SendInternalServerError(w)
		return
	}

	ret := redeemHandlerResponse{
		Success:               true,
		Error:                 "",
		SubscriptionExpiresAt: subscriptionEndUnix,
	}

	jRet, err := json.Marshal(&ret)
	if err != nil {
		util.HandleError(errors.Wrap(err, "error marshaling success json in redeemhandler"))
		util.SendInternalServerError(w)
		return
	}

	util.SendJSONResponse(w, jRet)
	return
}
