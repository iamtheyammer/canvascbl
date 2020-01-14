package canvas_tokens

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
	"time"
)

type CanvasToken struct {
	ID            uint64
	UserID        uint64
	GoogleUsersID uint64
	CanvasUserID  uint64
	Token         string
	ExpiresAt     time.Time
	InsertedAt    time.Time
}

type ListRequest struct {
	ID            uint64
	UserID        uint64
	GoogleUsersID uint64
	CanvasUserID  uint64
	Token         string

	Limit          uint64
	Offset         uint64
	NotExpiredOnly bool
	OrderBy        []string
}

// List lists Canvas tokens. Default order is insertion descending.
func List(db services.DB, req *ListRequest) (*[]CanvasToken, error) {
	q := util.Sq.
		Select(
			"id",
			"user_id",
			"google_users_id",
			"canvas_user_id",
			"token",
			"expires_at",
			"inserted_at",
		).
		From("canvas_tokens")

	if req.ID > 0 {
		q = q.Where(sq.Eq{"id": req.ID})
	}

	if req.UserID > 0 {
		q = q.Where(sq.Eq{"user_id": req.UserID})
	}

	if req.GoogleUsersID > 0 {
		q = q.Where(sq.Eq{"google_users_id": req.GoogleUsersID})
	}

	if req.CanvasUserID > 0 {
		q = q.Where(sq.Eq{"canvas_user_id": req.CanvasUserID})
	}

	if len(req.Token) > 0 {
		q = q.Where(sq.Eq{"token": req.Token})
	}

	if req.Limit > 0 {
		q = q.Limit(req.Limit)
	} else {
		req.Limit = services.DefaultSelectLimit
	}

	if req.Offset > 0 {
		q = q.Limit(req.Offset)
	}

	if req.NotExpiredOnly {
		q = q.Where("(expires_at IS NULL OR expires_at > NOW())")
	}

	if len(req.OrderBy) > 0 {
		q = q.OrderBy(req.OrderBy...)
	} else {
		q = q.OrderBy("inserted_at DESC")
	}

	query, args, err := q.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "error building list canvas tokens sql")
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "error opening list canvas tokens rows")
	}

	defer rows.Close()

	var cts []CanvasToken

	for rows.Next() {
		var (
			ct           CanvasToken
			userID       sql.NullInt64
			canvasUserID sql.NullInt64
			expiresAt    sql.NullTime
		)

		err := rows.Scan(
			&ct.ID,
			&userID,
			&ct.GoogleUsersID,
			&canvasUserID,
			&ct.Token,
			&expiresAt,
			&ct.InsertedAt,
		)
		if err != nil {
			return nil, errors.Wrap(err, "error scanning list canvas tokens rows")
		}

		if userID.Valid {
			ct.UserID = uint64(userID.Int64)
		}

		if canvasUserID.Valid {
			ct.CanvasUserID = uint64(canvasUserID.Int64)
		}

		if expiresAt.Valid {
			ct.ExpiresAt = expiresAt.Time
		}

		cts = append(cts, ct)
	}

	return &cts, nil
}
