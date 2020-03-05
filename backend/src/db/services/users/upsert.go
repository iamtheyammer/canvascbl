package users

import (
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
	InsertedAt time.Time
}

type UpsertObserveesRequest struct {
	Observees            []Observee
	ObserverCanvasUserID uint64
}

// UpsertProfile upserts a user profile
func UpsertProfile(db services.DB, ur *UpsertRequest) (*UpsertResponse, error) {
	query, args, err := util.Sq.
		Insert("users").
		SetMap(map[string]interface{}{
			"name":           ur.Name,
			"email":          ur.Email,
			"lti_user_id":    ur.LTIUserID,
			"canvas_user_id": ur.CanvasUserID,
		}).
		// normally would be ignore, but emails and names can change
		Suffix(
			"ON CONFLICT (lti_user_id) DO UPDATE SET name = EXCLUDED.name, email = EXCLUDED.email " +
				"RETURNING inserted_at",
		).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "error building upsert users sql")
	}

	row := db.QueryRow(query, args...)

	var resp UpsertResponse
	err = row.Scan(&resp.InsertedAt)
	if err != nil {
		return nil, errors.Wrap(err, "error executing upsert users sql")
	}

	return &resp, nil
}
