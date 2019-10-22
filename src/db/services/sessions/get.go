package sessions

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
	"time"
)

type Session struct {
	ID            uint64
	UserID        uint64
	SessionString string
	InsertedAt    time.Time
}

type GetRequest struct {
	ID            uint64
	SessionString string
}

// List gets a session by ID or SessionString. It returns nil, nil when no session is found for the params.
func List(db services.DB, req *GetRequest) (*Session, error) {
	if req.ID < 1 && len(req.SessionString) < 1 {
		return nil, nil
	}

	q := util.Sq.
		Select("id", "user_id", "session_string", "inserted_at").
		From("sessions")

	if req.ID > 0 {
		q = q.Where(sq.Eq{"id": req.ID})
	}

	if len(req.SessionString) > 0 {
		q = q.Where(sq.Eq{"session_string": req.SessionString})
	}

	query, args, err := q.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "error building get session sql")
	}

	row := db.QueryRow(query, args...)

	var sess Session
	err = row.Scan(&sess.ID, &sess.UserID, &sess.SessionString, &sess.InsertedAt)
	if err != nil {
		return nil, errors.Wrap(err, "error scanning session row")
	}

	return &sess, nil
}
