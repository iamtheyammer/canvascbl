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
	ID           uint64
	CanvasUserID uint64
	Token        string
	RefreshToken string
	ScopeVersion uint64
	ExpiresAt    time.Time
	InsertedAt   time.Time
}

type ListRequest struct {
	ID           uint64
	UserID       uint64
	CanvasUserID uint64
	Token        string
	RefreshToken string
	ScopeVersion uint64

	Limit      uint64
	Offset     uint64
	OrderBys   []string
	DistinctOn string
}

// List lists Canvas tokens. Default order is insertion descending.
func List(db services.DB, req *ListRequest) (*[]CanvasToken, error) {
	firstCol := "canvas_tokens.id"
	if len(req.DistinctOn) > 0 {
		firstCol = "DISTINCT ON (" + req.DistinctOn + ") " + firstCol
	}

	q := util.Sq.
		Select(
			firstCol,
			"canvas_tokens.canvas_user_id",
			"canvas_tokens.token",
			"canvas_tokens.refresh_token",
			"canvas_tokens.scope_version",
			"canvas_tokens.expires_at",
			"canvas_tokens.inserted_at",
		).
		From("canvas_tokens")

	if req.ID > 0 {
		q = q.Where(sq.Eq{"id": req.ID})
	}

	if req.UserID > 0 {
		q = q.Join("users ON canvas_tokens.canvas_user_id = users.canvas_user_id").
			Where(sq.Eq{"users.id": req.UserID})
	}

	if req.CanvasUserID > 0 {
		q = q.Where(sq.Eq{"canvas_user_id": req.CanvasUserID})
	}

	if len(req.Token) > 0 {
		q = q.Where(sq.Eq{"token": req.Token})
	}

	if len(req.RefreshToken) > 0 {
		q = q.Where(sq.Eq{"refresh_token": req.Token})
	}

	if req.ScopeVersion > 0 {
		q = q.Where(sq.Eq{"scope_version": req.ScopeVersion})
	}

	if req.Limit > 0 {
		q = q.Limit(req.Limit)
	} else {
		req.Limit = services.DefaultSelectLimit
	}

	if req.Offset > 0 {
		q = q.Limit(req.Offset)
	}

	if len(req.OrderBys) > 0 {
		q = q.OrderBy(req.OrderBys...)
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
			ct        CanvasToken
			expiresAt sql.NullTime
		)

		err := rows.Scan(
			&ct.ID,
			&ct.CanvasUserID,
			&ct.Token,
			&ct.RefreshToken,
			&ct.ScopeVersion,
			&expiresAt,
			&ct.InsertedAt,
		)
		if err != nil {
			return nil, errors.Wrap(err, "error scanning list canvas tokens rows")
		}

		if expiresAt.Valid {
			ct.ExpiresAt = expiresAt.Time
		}

		cts = append(cts, ct)
	}

	return &cts, nil
}
