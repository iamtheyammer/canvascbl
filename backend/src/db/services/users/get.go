package users

import (
	"database/sql"
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

type ListObserveesRequest struct {
	ID                   uint64
	ObserverCanvasUserID uint64
	ObserveeCanvasUserID uint64
	ObserveeName         string

	ActiveOnly bool
	Limit      uint64
	Offset     uint64
}

type User struct {
	ID   uint64
	Name string
	// almost always blank
	StripeCustomerID string
	Email            string
	LTIUserID        string
	CanvasUserID     uint64
	InsertedAt       time.Time
	Status           int
}

type Observee struct {
	ID             uint64
	ObserverUserID uint64
	CanvasUserID   uint64
	Name           string
	DeletedAt      time.Time
	InsertedAt     time.Time
}

func List(db services.DB, req *ListRequest) (*[]User, error) {
	q := util.Sq.
		Select("id", "name", "email", "lti_user_id", "canvas_user_id", "inserted_at", "status").
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

	var us []User

	for rows.Next() {
		var u User
		err = rows.Scan(
			&u.ID,
			&u.Name,
			&u.Email,
			&u.LTIUserID,
			&u.CanvasUserID,
			&u.InsertedAt,
			&u.Status,
		)
		if err != nil {
			return nil, errors.Wrap(err, "error scanning users")
		}

		us = append(us, u)
	}

	return &us, nil
}

func GetByStripeID(db services.DB, stripeID string) (*User, error) {
	query, args, err := util.Sq.
		Select("users.id AS user_id",
			"name",
			"stripe_customers.user_id AS stripe_user_id",
			"users.email",
			"users.lti_user_id",
			"users.canvas_user_id",
			"users.inserted_at",
			"status").
		From("users").
		Join("stripe_customers ON users.id = stripe_customers.user_id").
		Where(sq.Eq{"stripe_customers.stripe_id": stripeID}).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "error building get user by stripe id sql")
	}

	row := db.QueryRow(query, args...)

	var u User
	err = row.Scan(
		&u.ID,
		&u.Name,
		&u.StripeCustomerID,
		&u.Email,
		&u.LTIUserID,
		&u.CanvasUserID,
		&u.InsertedAt,
		&u.Status,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrap(err, "error scanning get user by stripe id")
	}

	return &u, nil
}

func ListObservees(db services.DB, req *ListObserveesRequest) (*[]Observee, error) {
	q := util.Sq.
		Select(
			"id",
			"observer_canvas_user_id",
			"observee_canvas_user_id",
			"observee_name",
			"deleted_at",
			"inserted_at",
		).
		From("observees")

	if req.ID != 0 {
		q = q.Where(sq.Eq{"id": req.ID})
	}

	if req.ObserverCanvasUserID != 0 {
		q = q.Where(sq.Eq{"observer_canvas_user_id": req.ObserverCanvasUserID})
	}

	if req.ObserveeCanvasUserID != 0 {
		q = q.Where(sq.Eq{"observee_canvas_user_id": req.ObserverCanvasUserID})
	}

	if req.ActiveOnly {
		q = q.Where(sq.Eq{"deleted_at": nil})
	}

	if req.Limit != 0 {
		q = q.Limit(req.Limit)
	}

	if req.Offset != 0 {
		q = q.Offset(req.Offset)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "error building list observees sql")
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "error executing list observees sql")
	}

	var os []Observee

	for rows.Next() {
		var (
			o         Observee
			name      sql.NullString
			deletedAt sql.NullTime
		)

		err := rows.Scan(
			&o.ID,
			&o.ObserverUserID,
			&o.CanvasUserID,
			&name,
			&deletedAt,
			&o.InsertedAt,
		)
		if err != nil {
			return nil, errors.Wrap(err, "error scanning list observees sql")
		}

		if name.Valid {
			o.Name = name.String
		}

		if deletedAt.Valid {
			o.DeletedAt = deletedAt.Time
		}

		os = append(os, o)
	}

	return &os, nil
}
