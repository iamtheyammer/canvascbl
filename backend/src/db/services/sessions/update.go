package sessions

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
)

type UpdateRequest struct {
	CanvasUserID  uint64
	GoogleUsersID uint64

	WhereSessionString string
	WhereID            uint64
}

func Update(db services.DB, req *UpdateRequest) error {
	q := util.Sq.
		Update("sessions")

	if req.WhereID > 0 {
		q = q.Where(sq.Eq{"id": req.WhereID})
	}

	if len(req.WhereSessionString) > 0 {
		q = q.Where(sq.Eq{"session_string": req.WhereSessionString})
	}

	if req.CanvasUserID > 0 {
		q = q.Set("canvas_user_id", req.CanvasUserID)
	}

	if req.GoogleUsersID > 0 {
		q = q.Set("google_users_id", req.GoogleUsersID)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return errors.Wrap(err, "error building update session sql")
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "error executing update session sql")
	}

	return nil
}
