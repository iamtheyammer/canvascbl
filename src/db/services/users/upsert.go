package users

import (
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/orautil"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
)

type UpsertRequest struct {
	Name         string
	Email        string
	LTIUserID    string
	CanvasUserID int64
}

type UpsertObserveesRequest struct {
	Observees            []Observee
	ObserverCanvasUserID uint64
}

// UpsertProfile upserts a user profile
func UpsertProfile(db services.DB, req *UpsertRequest) error {
	query, args, err := orautil.BuildUpsert(
		"users",
		[]orautil.UpsertableColumn{
			{
				Name:  "canvas_user_id",
				Value: req.CanvasUserID,
			},
		},
		[]orautil.UpsertableColumn{
			{
				"name",
				req.Name,
				false,
			},
			{
				"email",
				req.Email,
				false,
			},
			{
				"lti_user_id",
				req.LTIUserID,
				true,
			},
			{
				"canvas_user_id",
				req.CanvasUserID,
				true,
			},
		},
	)
	if err != nil {
		return errors.Wrap(err, "error building upsert users sql")
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "error executing upsert users sql")
	}

	return nil
}

// UpsertUserObservees upserts a user's observees.
func UpsertUserObservees(db services.DB, req *UpsertObserveesRequest) error {
	if len(req.Observees) < 1 {
		return nil
	}

	q := util.Sq.
		Insert("observees").
		Columns("observer_canvas_user_id", "observee_canvas_user_id", "observee_name").
		Suffix("ON CONFLICT ON CONSTRAINT observees_observer_canvas_user_id_observee_canvas_user_id_key " +
			"DO UPDATE SET observee_name = excluded.observee_name")

	for _, o := range req.Observees {
		q = q.Values(req.ObserverCanvasUserID, o.CanvasUserID, o.Name)
	}

	checks := []orautil.UpsertableColumn{
		{
			Name: "observer_canvas_user_id",

		}
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
