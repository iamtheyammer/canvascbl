package subscriptions

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
	"time"
)

type Subscription struct {
	ID                 uint64
	StripeID           string
	UserID             uint64
	CustomerStripeID   string
	Plan               string
	Status             string
	Price              uint64
	CurrentPeriodStart time.Time
	CurrentPeriodEnd   time.Time
	TrialEnd           time.Time
	CanceledAt         time.Time
	InsertedAt         time.Time
	UpdatedAt          time.Time
}

type GetRequest struct {
	ID               uint64
	StripeID         string
	UserID           uint64
	CustomerStripeID string
	Status           string
	ActiveOnly       bool
	Limit            uint64
	Offset           uint64
	OrderBy          string
}

func Get(db services.DB, req *GetRequest) (*[]Subscription, error) {
	q := util.Sq.
		Select(
			"id",
			"stripe_id",
			"user_id",
			"customer_stripe_id",
			"plan",
			"status",
			"price",
			"current_period_start",
			"current_period_end",
			"trial_end",
			"canceled_at",
			"inserted_at",
			"updated_at",
		).
		From("subscriptions").
		Limit(services.DefaultSelectLimit)

	if req.ID > 0 {
		q = q.Where(sq.Eq{"id": req.ID})
	}

	if len(req.StripeID) > 0 {
		q = q.Where(sq.Eq{"stripe_id": req.StripeID})
	}

	if req.UserID > 0 {
		q = q.Where(sq.Eq{"user_id": req.UserID})
	}

	if len(req.CustomerStripeID) > 0 {
		q = q.Where(sq.Eq{"customer_stripe_id": req.CustomerStripeID})
	}

	if len(req.Status) > 0 {
		q = q.Where(sq.Eq{"status": req.Status})
	}

	if req.ActiveOnly {
		q = q.Where(sq.Eq{"status": []string{"active", "trialing"}}).
			Where("current_period_end > NOW()")
	}

	if req.Limit > 0 {
		q = q.Limit(req.Limit)
	}

	if req.Offset > 0 {
		q = q.Offset(req.Offset)
	}

	if len(req.OrderBy) > 0 {
		q = q.OrderBy(req.OrderBy)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "error building get subscriptions sql")
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "error executing get subscriptions sql")
	}

	defer rows.Close()

	var ss []Subscription

	for rows.Next() {
		var (
			s                    Subscription
			trialEnd, canceledAt sql.NullTime
			price                sql.NullInt64
		)

		err := rows.Scan(
			&s.ID,
			&s.StripeID,
			&s.UserID,
			&s.CustomerStripeID,
			&s.Plan,
			&s.Status,
			&price,
			&s.CurrentPeriodStart,
			&s.CurrentPeriodEnd,
			&trialEnd,
			&canceledAt,
			&s.InsertedAt,
			&s.UpdatedAt,
		)

		if err != nil {
			return nil, errors.Wrap(err, "error scanning subscriptions")
		}

		if trialEnd.Valid {
			s.TrialEnd = trialEnd.Time
		}

		if canceledAt.Valid {
			s.CanceledAt = canceledAt.Time
		}

		if price.Valid {
			s.Price = uint64(price.Int64)
		}

		ss = append(ss, s)
	}

	return &ss, nil
}
