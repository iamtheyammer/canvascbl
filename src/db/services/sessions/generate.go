package sessions

import (
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

type GenerateRequest struct {
	CanvasUserID  uint64
	GoogleUsersID uint64
}

// Generate generates a session for the given userID using the given database
func Generate(db services.DB, req *GenerateRequest) (*string, error) {
	u := uuid.NewV4().String()

	q := util.Sq.
		Insert("sessions")

	if req.CanvasUserID > 0 && req.GoogleUsersID > 0 {
		q = q.SetMap(map[string]interface{}{
			"canvas_user_id":  req.CanvasUserID,
			"google_users_id": req.GoogleUsersID,
			"session_string":  u,
		})
	} else if req.GoogleUsersID > 0 {
		q = q.SetMap(map[string]interface{}{
			"google_users_id": req.GoogleUsersID,
			"session_string":  u,
		})
	} else if req.CanvasUserID > 0 {
		q = q.SetMap(map[string]interface{}{
			"canvas_user_id": req.CanvasUserID,
			"session_string": u,
		})
	} else {
		return nil, nil
	}

	query, args, err := q.ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "error building insert session sql")
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "error executing insert session sql")
	}

	return &u, nil
}
