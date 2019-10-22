package sessions

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
)

type VerifiedSession struct {
	UserID               uint64
	CanvasUserID         uint64
	Email                string
	HasValidSubscription bool
	SubscriptionStatus   string
}

func Verify(db services.DB, sessionString string) (*VerifiedSession, error) {
	if len(sessionString) != 36 {
		return nil, nil
	}

	query, args, err := util.Sq.
		Select(
			"users.id AS user_id",
			"users.canvas_user_id AS canvas_user_id",
			"users.email AS email",
			"(CASE WHEN subscriptions.status IN ('active', 'trialing') "+
				"AND subscriptions.current_period_end > NOW() THEN "+
				"TRUE ELSE FALSE END) AS has_valid_subscription",
			"subscriptions.status AS subscription_status",
		).
		From("subscriptions").
		Join("users ON subscriptions.user_id = users.id").
		Join("sessions ON users.canvas_user_id = sessions.user_id").
		Where(sq.Eq{"sessions.session_string": sessionString}).
		OrderBy("has_valid_subscription DESC").
		Limit(1).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "error building verify session sql")
	}

	row := db.QueryRow(query, args...)

	var (
		vs                 VerifiedSession
		subscriptionStatus sql.NullString
	)

	err = row.Scan(
		&vs.UserID,
		&vs.CanvasUserID,
		&vs.Email,
		&vs.HasValidSubscription,
		&subscriptionStatus,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrap(err, "error scanning verify session row sql")
	}

	if subscriptionStatus.Valid {
		vs.SubscriptionStatus = subscriptionStatus.String
	}

	return &vs, nil
}
