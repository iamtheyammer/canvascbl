package notifications

import (
	"fmt"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
)

// InsertNotificationSettingsRequest is the request for InsertNotificationSettings.
type InsertNotificationSettingsRequest struct {
	UserID uint64
	TypeID uint64
	Medium Medium
}

// InsertNotficationSettings inserts notification settings.
func InsertNotificationSettings(db services.DB, req *InsertNotificationSettingsRequest) error {
	query, args, err := util.Sq.
		Insert("notification_settings").
		SetMap(map[string]interface{}{
			"user_id":              req.UserID,
			"notification_type_id": req.TypeID,
			"medium":               req.Medium,
		}).
		Suffix("ON CONFLICT DO NOTHING").
		ToSql()
	if err != nil {
		return fmt.Errorf("error building insert notification settings sql: %w", err)
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("error executing insert notification settings sql: %w", err)
	}

	return nil
}
