package users

import (
	"fmt"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
	"time"
)

type UpsertRequest struct {
	Name         string
	Email        string
	LTIUserID    string
	CanvasUserID int64
}

// UpsertResponse contains some user data sometimes needed after an upsert.
type UpsertResponse struct {
	UserID     uint64
	InsertedAt time.Time
}

type UpsertObserveesRequest struct {
	Observees            []Observee
	ObserverCanvasUserID uint64
}

// UpsertProfile wraps UpsertMultipleProfiles for one user only
func UpsertProfile(db services.DB, ur *UpsertRequest, returnInsertedAt bool) (*UpsertResponse, error) {
	resp, err := UpsertMultipleProfiles(db, &[]UpsertRequest{*ur}, returnInsertedAt)
	if err != nil {
		return nil, err
	}

	if resp != nil && len(*resp) > 0 {
		return &(*resp)[0], nil
	} else {
		return nil, nil
	}
}

// UpsertMultipleProfiles upserts multiple user profiles
func UpsertMultipleProfiles(db services.DB, ur *[]UpsertRequest, returnInsertedAt bool) (*[]UpsertResponse, error) {
	suffix := "ON CONFLICT (canvas_user_id) DO UPDATE SET name = EXCLUDED.name, email = EXCLUDED.email"
	if returnInsertedAt {
		suffix += " RETURNING id, inserted_at"
	}

	q := util.Sq.
		Insert("users").
		Columns(
			"name",
			"email",
			"lti_user_id",
			"canvas_user_id",
		).
		Suffix(suffix)

	for _, r := range *ur {
		var LTIUserID interface{}
		if len(r.LTIUserID) > 0 {
			LTIUserID = r.LTIUserID
		}

		q = q.Values(r.Name, r.Email, LTIUserID, r.CanvasUserID)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building upsert users sql: %w", err)
	}

	if returnInsertedAt {
		rows, err := db.Query(query, args...)
		if err != nil {
			return nil, fmt.Errorf("error executing upsert users (returning) sql: %w", err)
		}

		var resp []UpsertResponse
		for rows.Next() {
			var r UpsertResponse

			err = rows.Scan(&r.UserID, &r.InsertedAt)
			if err != nil {
				return nil, fmt.Errorf("error scanning upsert users sql: %w", err)
			}

			resp = append(resp, r)
		}

		return &resp, nil
	} else {
		_, err = db.Exec(query, args...)
		if err != nil {
			return nil, fmt.Errorf("error executing upsert users sql: %w", err)
		}

		return nil, nil
	}
}

// UpsertUserObservees inserts a user's observees.
func UpsertUserObservees(db services.DB, req *UpsertObserveesRequest) error {
	if len(req.Observees) < 1 {
		return nil
	}

	q := util.Sq.
		Insert("observees").
		Columns("observer_canvas_user_id", "observee_canvas_user_id", "observee_name").
		Suffix("ON CONFLICT (observer_canvas_user_id, observee_canvas_user_id) DO UPDATE SET " +
			"observee_canvas_user_id = EXCLUDED.observee_canvas_user_id, " +
			"observee_name = EXCLUDED.observee_name")

	for _, o := range req.Observees {
		q = q.Values(req.ObserverCanvasUserID, o.CanvasUserID, o.Name)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return errors.Wrap(err, "error building insert observees sql")
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "error executing insert observees sql")
	}

	return nil
}
