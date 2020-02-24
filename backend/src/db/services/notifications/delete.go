package notifications

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
)

// DeleteNotificationSettingRequest is the request for DeleteNotificationSetting.
type DeleteNotificationSettingRequest struct {
	UserID uint64
	TypeID uint64
	Medium Medium
}

// DeleteNotificationSetting deletes a notification setting.
func DeleteNotificationSetting(db services.DB, req *DeleteNotificationSettingRequest) error {
	q := util.Sq.
		Delete("notification_settings")

	if req.UserID > 0 {
		q = q.Where(sq.Eq{"user_id": req.UserID})
	}

	if req.TypeID > 0 {
		q = q.Where(sq.Eq{"notification_type_id": req.TypeID})
	}

	if req.Medium != zeroMedium {
		q = q.Where(sq.Eq{"medium": req.Medium})
	}

	query, args, err := q.ToSql()
	if err != nil {
		return fmt.Errorf("error building delete notification setting sql: %w", err)
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("error executing delete notification setting sql: %w", err)
	}

	return nil
}
