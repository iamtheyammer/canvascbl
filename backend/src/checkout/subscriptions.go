package checkout

import (
	"encoding/json"
	"github.com/iamtheyammer/canvascbl/backend/src/db"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/subscriptions"
	"github.com/iamtheyammer/canvascbl/backend/src/middlewares"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"github.com/stripe/stripe-go/sub"
	"net/http"
)

type subscription struct {
	ID                 uint64      `json:"id"`
	StripeID           string      `json:"stripeId"`
	UserID             uint64      `json:"userId"`
	CustomerStripeID   string      `json:"customerStripeId"`
	Plan               string      `json:"plan"`
	Status             string      `json:"status"`
	Price              interface{} `json:"price"`
	CurrentPeriodStart int64       `json:"currentPeriodStart"`
	CurrentPeriodEnd   int64       `json:"currentPeriodEnd"`
	TrialEnd           interface{} `json:"trialEnd"`
	CanceledAt         interface{} `json:"canceledAt"`
	InsertedAt         int64       `json:"insertedAt"`
	UpdatedAt          int64       `json:"updatedAt"`
}

func CancelSubscriptionHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	session := middlewares.Session(w, req)
	if session == nil {
		return
	}

	subs, err := db.GetSubscriptions(&subscriptions.GetRequest{
		UserID:     session.UserID,
		ActiveOnly: true,
	})
	if err != nil {
		util.HandleError(errors.Wrap(err, "error getting subscriptions"))
		util.SendInternalServerError(w)
		return
	}

	if len(*subs) < 1 {
		util.SendNotFoundWithReason(w, "you don't have an active subscription")
		return
	}

	_, err = sub.Cancel((*subs)[0].StripeID, nil)
	if err != nil {
		util.HandleError(errors.Wrap(err, "error cancelling subscription"))
		util.SendInternalServerError(w)
		return
	}

	util.SendNoContent(w)
	return
}

func ListSubscriptionsHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	session := middlewares.Session(w, req)
	if session == nil {
		return
	}

	dbSubs, err := db.GetSubscriptions(&subscriptions.GetRequest{
		UserID:     session.UserID,
		ActiveOnly: true,
	})
	if err != nil {
		util.HandleError(errors.Wrap(err, "error listing subscriptions"))
		util.SendInternalServerError(w)
		return
	}
	dbs := *dbSubs

	if len(dbs) < 1 {
		util.SendJSONResponse(w, []byte("[]"))
		return
	}

	var subs []subscription

	for _, s := range dbs {
		su := subscription{
			ID:                 s.ID,
			StripeID:           s.StripeID,
			UserID:             s.UserID,
			CustomerStripeID:   s.CustomerStripeID,
			Plan:               s.Plan,
			Status:             s.Status,
			Price:              s.Price,
			CurrentPeriodStart: s.CurrentPeriodStart.Unix(),
			CurrentPeriodEnd:   s.CurrentPeriodEnd.Unix(),
			TrialEnd:           s.TrialEnd.Unix(),
			CanceledAt:         s.CanceledAt.Unix(),
			InsertedAt:         s.InsertedAt.Unix(),
			UpdatedAt:          s.UpdatedAt.Unix(),
		}

		if s.Price < 0 {
			su.Price = nil
		}

		if s.TrialEnd.Unix() < 0 {
			su.TrialEnd = nil
		}

		if s.CanceledAt.Unix() < 0 {
			su.CanceledAt = nil
		}

		subs = append(subs, su)
	}

	jsubs, err := json.Marshal(subs)
	if err != nil {
		util.HandleError(errors.Wrap(err, "error marshaling subscriptions"))
		util.SendInternalServerError(w)
		return
	}

	util.SendJSONResponse(w, jsubs)
	return
}
