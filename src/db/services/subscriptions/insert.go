package subscriptions

import (
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
	"time"
)

type InsertRequest struct {
	StripeID           string
	UserID             uint64
	CustomerStripeID   string
	Plan               string
	Status             string
	Price              uint64
	CurrentPeriodStart uint64
	CurrentPeriodEnd   uint64
	TrialEnd           uint64
	CanceledAt         uint64
}

func Insert(db services.DB, req *InsertRequest) error {
	var price interface{}
	if req.Price > 0 {
		price = req.Price
	} else {
		price = nil
	}

	var trialEnd interface{}
	if req.TrialEnd > 0 {
		trialEnd = time.Unix(int64(req.TrialEnd), 0)
	} else {
		trialEnd = nil
	}

	var canceledAt interface{}
	if req.CanceledAt > 0 {
		canceledAt = time.Unix(int64(req.CanceledAt), 0)
	} else {
		canceledAt = nil
	}

	query, args, err := util.Sq.
		Insert("subscriptions").
		SetMap(map[string]interface{}{
			"stripe_id":            req.StripeID,
			"user_id":              req.UserID,
			"customer_stripe_id":   req.CustomerStripeID,
			"plan":                 req.Plan,
			"status":               req.Status,
			"price":                price,
			"current_period_start": time.Unix(int64(req.CurrentPeriodStart), 0),
			"current_period_end":   time.Unix(int64(req.CurrentPeriodEnd), 0),
			"trial_end":            trialEnd,
			"canceled_at":          canceledAt,
		}).ToSql()
	if err != nil {
		return errors.Wrap(err, "error building insert subscription sql")
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "error executing insert subscription sql")
	}

	return nil
}
