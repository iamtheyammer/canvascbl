package users

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
	"time"
)

type ListRequest struct {
	ID           uint64
	Email        string
	LTIUserID    string
	CanvasUserID uint64
	Limit        uint64
	Offset       uint64
}

type User struct {
	ID           uint64
	Name         string
	Email        string
	LTIUserID    string
	CanvasUserID uint64
	InsertedAt   time.Time
}

func List(db services.DB, req *ListRequest) (*[]User, error) {
	q := util.Sq.
		Select("id", "name", "email", "lti_user_id", "canvas_user_id", "inserted_at").
		From("users").
		Limit(services.DefaultSelectLimit)

	if req.ID > 0 {
		q = q.Where(sq.Eq{"id": req.ID})
	}

	if len(req.Email) > 0 {
		q = q.Where(sq.Eq{"email": req.Email})
	}

	if len(req.LTIUserID) > 0 {
		q = q.Where(sq.Eq{"lti_user_id": req.LTIUserID})
	}

	if req.CanvasUserID > 0 {
		q = q.Where(sq.Eq{"canvas_user_id": req.CanvasUserID})
	}

	if req.Limit > 0 {
		q = q.Limit(req.Limit)
	}

	if req.Offset > 0 {
		q = q.Offset(req.Offset)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "error building list users sql")
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "error querying db")
	}

	defer rows.Close()

	var users []User

	for rows.Next() {
		var u User
		err = rows.Scan(
			&u.ID,
			&u.Name,
			&u.Email,
			&u.LTIUserID,
			&u.CanvasUserID,
			&u.InsertedAt,
		)
		if err != nil {
			return nil, errors.Wrap(err, "error scanning users")
		}

		users = append(users, u)
	}

	return &users, nil
}
