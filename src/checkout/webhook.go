package checkout

import (
	"encoding/json"
	"github.com/iamtheyammer/canvascbl/backend/src/db"
	"github.com/iamtheyammer/canvascbl/backend/src/env"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/webhook"
	"io/ioutil"
	"net/http"
)

const StripeWebhookMaxBodyBytes = int64(65536)

func StripeWebhookHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	req.Body = http.MaxBytesReader(w, req.Body, StripeWebhookMaxBodyBytes)
	payload, err := ioutil.ReadAll(req.Body)
	if err != nil {
		util.HandleError(errors.Wrap(err, "error reading request body from stripe webhook"))
		util.SendInternalServerError(w)
		return
	}

	event, err := webhook.ConstructEvent(payload, req.Header.Get("Stripe-Signature"), env.StripeWebhookSecret)
	if err != nil {
		util.HandleError(errors.Wrap(err, "error verifying stripe signature on webhook pull"))
		util.SendBadRequest(w, "error verifying Stripe-Signature")
		return
	}

	if err := json.Unmarshal(payload, &event); err != nil {
		util.HandleError(errors.Wrap(err, "error unmarshaling event from stripe webhook"))
		util.SendBadRequest(w, "malformed payload")
		return
	}

	switch event.Type {
	case "customer.subscription.created":
		sub, err := stripeWebhookProcessSubscription(event.Data.Raw)
		if err != nil {
			util.HandleError(errors.Wrap(err, "error unmarshaling customer.subscription.created event from stripe webhook"))
			util.SendBadRequest(w, "unable to unmarshal into stripe.Subscription")
			return
		}
		err = db.CheckoutWebhookInsertSubscription(*sub)
		if err != nil {
			util.HandleError(errors.Wrap(err, "error inserting subscription"))
			return
		}
		return
	case "customer.subscription.updated":
		sub, err := stripeWebhookProcessSubscription(event.Data.Raw)
		if err != nil {
			util.HandleError(errors.Wrap(err, "error unmarshaling customer.subscription.updated event from stripe webhook"))
			util.SendBadRequest(w, "unable to unmarshal into stripe.Subscription")
			return
		}
		err = db.CheckoutWebhookUpdateSubscription(*sub)
		if err != nil {
			util.HandleError(errors.Wrap(err, "error updating subscription"))
			return
		}
		return
	case "customer.subscription.deleted":
		sub, err := stripeWebhookProcessSubscription(event.Data.Raw)
		if err != nil {
			util.HandleError(errors.Wrap(err, "error unmarshaling customer.subscription.deleted event from stripe webhook"))
			util.SendBadRequest(w, "unable to unmarshal into stripe.Subscription")
			return
		}
		err = db.CheckoutWebhookUpdateSubscription(*sub)
		if err != nil {
			util.HandleError(errors.Wrap(err, "error deleting subscription"))
			return
		}
		return
	default:
		util.SendBadRequest(w, "that event isn't handled by this endpoint at this time")
		return
	}
}

func stripeWebhookProcessSubscription(msg json.RawMessage) (*stripe.Subscription, error) {
	var sub stripe.Subscription
	err := json.Unmarshal(msg, &sub)
	if err != nil {
		return nil, errors.Wrap(err, "error processing (unmarshaling) stripe subscription")
	}

	return &sub, nil
}
