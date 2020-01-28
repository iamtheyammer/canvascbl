package users

import (
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
)

// InsertUserObservees inserts a user's observees.
func InsertUserObservees(db services.DB, req *UpsertObserveesRequest) error {
	if len(req.Observees) < 1 {
		return nil
	}

	q := util.Sq.
		Insert("observees").
		Columns("observer_canvas_user_id", "observee_canvas_user_id", "observee_name").
		Suffix("ON CONFLICT DO NOTHING")

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
