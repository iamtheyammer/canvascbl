package stripe_customers

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
	"time"
)

type StripeCustomer struct {
	ID       uint64
	StripeID string
	// users.id
	UserID     string
	InsertedAt time.Time
}

// Get takes a user ID (users.id) and returns a StripeCustomer, or nil if that user doesn't have a stripe customer.
// If the user has more than one customer, it returns the newest one.
func Get(db services.DB, userID uint64) (*StripeCustomer, error) {
	query, args, err := util.Sq.
		Select("id", "stripe_id", "user_id", "inserted_at").
		From("stripe_customers").
		Where(sq.Eq{"user_id": userID}).
		OrderBy("inserted_at DESC").
		Limit(1).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "error building get stripe customer by user id sql")
	}

	row := db.QueryRow(query, args...)

	var cust StripeCustomer

	err = row.Scan(
		&cust.ID,
		&cust.StripeID,
		&cust.UserID,
		&cust.InsertedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrap(err, "error scanning get stripe customer by user id sql")
	}

	return &cust, nil
}
