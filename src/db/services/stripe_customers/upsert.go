package stripe_customers

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
)

type UpsertRequest struct {
	// ID of the customer from Stripe
	StripeID string
	// users.id-- can be provided OR an email-- NEVER BOTH
	UserID uint64
	// user's email-- can be provided OR UserID-- NEVER BOTH
	Email string
}

// Upsert upserts a customer, but does nothing if there is a duplicate, as user IDs and customer IDs don't change.
func Upsert(db services.DB, req *UpsertRequest) error {
	if len(req.Email) < 1 && len(req.StripeID) < 1 {
		return errors.New("neither email or stripeID were passed in")
	}

	q := util.Sq.
		Insert("stripe_customers").
		Suffix("ON CONFLICT DO NOTHING").
		Columns("stripe_id", "user_id")

	if req.UserID > 0 {
		q = q.Values(req.StripeID, req.UserID)
	}

	if len(req.Email) > 0 {
		q = q.Values(req.StripeID, sq.Expr("(SELECT id FROM users WHERE email = ?)", req.Email))
	}

	query, args, err := q.ToSql()
	if err != nil {
		return errors.Wrap(err, "error building upsert customers sql")
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "error executing upsert customers sql")
	}

	return nil
}
