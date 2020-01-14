package canvas_tokens

import (
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
	"time"
)

type InsertRequest struct {
	UserID *uint64
	// google_users.id
	GoogleUsersID uint64
	CanvasUserID  *uint64
	Token         string
	ExpiresAt     *time.Time
}

func Insert(db services.DB, req *InsertRequest) error {
	query, args, err := util.Sq.
		Insert("canvas_tokens").
		SetMap(map[string]interface{}{
			"user_id":         req.UserID,
			"google_users_id": req.GoogleUsersID,
			"canvas_user_id":  req.CanvasUserID,
			"token":           req.Token,
			"expires_at":      req.ExpiresAt,
		}).
		Suffix("ON CONFLICT DO NOTHING").
		ToSql()
	if err != nil {
		return errors.Wrap(err, "error building insert canvas token sql")
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "error executing insert canvas token sql")
	}

	return nil
}
