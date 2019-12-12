package sessions

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
)

type VerifiedSession struct {
	SessionString string
	// users.id
	UserID               uint64
	CanvasUserID         uint64
	UserStatus           int
	GoogleUsersID        string
	Email                string
	HasValidSubscription bool
	SubscriptionStatus   string
	SessionIsExpired     bool
}

func Verify(db services.DB, sessionString string) (*VerifiedSession, error) {
	if len(sessionString) != 36 {
		return nil, nil
	}

	query, args, err := util.Sq.
		Select(
			"users.id AS user_id",
			"users.canvas_user_id AS canvas_user_id",
			"users.status AS user_status",
			"google_users.id AS google_users_id",
			"users.email AS email",
			"(CASE WHEN subscriptions.status IN ('active', 'trialing') "+
				"AND subscriptions.current_period_end > NOW() THEN "+
				"TRUE ELSE FALSE END) AS has_valid_subscription",
			"subscriptions.status AS subscription_status",
			"(CASE WHEN sessions.inserted_at + interval '2 weeks' < NOW() THEN TRUE ELSE FALSE END) "+
				"AS session_is_expired",
		).
		From("subscriptions").
		RightJoin("users ON subscriptions.user_id = users.id").
		RightJoin("google_users ON LOWER(users.email) = LOWER(google_users.email)").
		Join("sessions ON (users.canvas_user_id = sessions.canvas_user_id OR google_users.id = sessions.google_users_id)").
		Where(sq.Eq{"sessions.session_string": sessionString}).
		OrderBy("has_valid_subscription DESC").
		Limit(1).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "error building verify session sql")
	}

	row := db.QueryRow(query, args...)

	var (
		vs                                       VerifiedSession
		userID, canvasUserID                     sql.NullInt64
		email, subscriptionStatus, googleUsersID sql.NullString
	)

	err = row.Scan(
		&userID,
		&canvasUserID,
		&vs.UserStatus,
		&googleUsersID,
		&email,
		&vs.HasValidSubscription,
		&subscriptionStatus,
		&vs.SessionIsExpired,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrap(err, "error scanning verify session row sql")
	}

	if userID.Valid {
		vs.UserID = uint64(userID.Int64)
	}

	if canvasUserID.Valid {
		vs.CanvasUserID = uint64(canvasUserID.Int64)
	}

	if googleUsersID.Valid {
		vs.GoogleUsersID = googleUsersID.String
	}

	if email.Valid {
		vs.Email = email.String
	}

	if subscriptionStatus.Valid {
		vs.SubscriptionStatus = subscriptionStatus.String
	}

	vs.SessionString = sessionString

	return &vs, nil
}
