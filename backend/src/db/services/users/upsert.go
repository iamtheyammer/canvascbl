package users

import (
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
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
func UpsertProfile(db services.DB, ur *UpsertRequest) error {
	query, args, err := util.Sq.
		Insert("users").
		SetMap(map[string]interface{}{
			"name":           ur.Name,
			"email":          ur.Email,
			"lti_user_id":    ur.LTIUserID,
			"canvas_user_id": ur.CanvasUserID,
		}).
		// normally would be ignore, but emails and names can change
		Suffix("ON CONFLICT ON CONSTRAINT users_lti_user_id_key DO UPDATE SET name = ?, email = ?", ur.Name, ur.Email).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "error building upsert users sql")
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "error executing upsert users sql")
	}

	return nil
}
