package sessions

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
)

type VerifiedSession struct {
	SessionString string `json:"-"`
	// users.id
	UserID               uint64 `json:"user_id"`
	CanvasUserID         uint64 `json:"canvas_user_id"`
	UserStatus           int    `json:"status"`
	Email                string `json:"email"`
	HasValidSubscription bool   `json:"has_valid_subscription"`
	SessionIsExpired     bool   `json:"-"`
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
			//"google_users.id AS google_users_id",
			"users.email AS email",
			"users.has_valid_subscription AS has_valid_subscription",
			"(CASE WHEN sessions.inserted_at + interval '2 weeks' < NOW() THEN TRUE ELSE FALSE END) "+
				"AS session_is_expired",
		).
		From("users").
		//RightJoin("google_users ON LOWER(users.email) = LOWER(google_users.email)").
		//Join("sessions ON (users.canvas_user_id = sessions.canvas_user_id OR google_users.id = sessions.google_users_id)").
		Join("sessions ON users.canvas_user_id = sessions.canvas_user_id").
		Where(sq.Eq{"sessions.session_string": sessionString}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "error building verify session sql")
	}

	row := db.QueryRow(query, args...)

	var (
		vs                   VerifiedSession
		userID, canvasUserID sql.NullInt64
		email                sql.NullString
		//googleUsersID             sql.NullString
	)

	err = row.Scan(
		&userID,
		&canvasUserID,
		&vs.UserStatus,
		//&googleUsersID,
		&email,
		&vs.HasValidSubscription,
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

	//if googleUsersID.Valid {
	//	vs.GoogleUsersID = googleUsersID.String
	//}

	if email.Valid {
		vs.Email = email.String
	}

	vs.SessionString = sessionString

	// COVID-19
	vs.HasValidSubscription = true

	return &vs, nil
}
