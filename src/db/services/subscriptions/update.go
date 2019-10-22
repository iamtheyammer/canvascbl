package subscriptions

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
	"time"
)

type UpdateRequest struct {
	Identifier struct {
		StripeID string
		ID       uint64
	}
	Values struct {
		UserID             uint64
		Status             string
		Price              uint64
		CurrentPeriodStart uint64
		CurrentPeriodEnd   uint64
		TrialEnd           uint64
		CanceledAt         uint64
	}
}

// Update will update all mutable fields
func Update(db services.DB, req *UpdateRequest) error {
	if req.Identifier.ID < 1 && len(req.Identifier.StripeID) < 1 {
		return errors.New("no identifier passed in")
	}

	q := util.Sq.
		Update("subscriptions")

	if req.Identifier.ID > 0 {
		q = q.Where(sq.Eq{"id": req.Identifier.ID})
	}

	if len(req.Identifier.StripeID) > 1 {
		q = q.Where(sq.Eq{"stripe_id": req.Identifier.StripeID})
	}

	if req.Values.UserID > 1 {
		q = q.Set("user_id", req.Values.UserID)
	}

	if len(req.Values.Status) > 1 {
		q = q.Set("status", req.Values.Status)
	}

	if req.Values.Price > 0 {
		q = q.Set("price", req.Values.Price)
	}

	if req.Values.CurrentPeriodStart > 0 {
		q = q.Set("current_period_start", time.Unix(int64(req.Values.CurrentPeriodStart), 0))
	}

	if req.Values.CurrentPeriodEnd > 0 {
		q = q.Set("current_period_end", time.Unix(int64(req.Values.CurrentPeriodEnd), 0))
	}

	if req.Values.TrialEnd > 0 {
		q = q.Set("trial_end", time.Unix(int64(req.Values.TrialEnd), 0))
	}

	if req.Values.CanceledAt > 0 {
		q = q.Set("canceled_at", time.Unix(int64(req.Values.CanceledAt), 0))
	}

	query, args, err := q.ToSql()
	if err != nil {
		return errors.Wrap(err, "error building update subscription sql")
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "error executing update subscription sql")
	}

	return nil
}
