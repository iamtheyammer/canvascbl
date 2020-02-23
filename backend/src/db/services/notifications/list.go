package notifications

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"time"
)

// Type represents a notification type.
type Type uint64

// Medium represents a notification medium.
type Medium string

const (
	// zeroType is Type's zero value
	zeroType = Type(0)
	// GradeChangeNotificationID is the ID for the grade_change notification type.
	TypeGradeChange = Type(1)

	// zeroMedium is Medium's zero value
	zeroMedium = Medium("")
	// MediumEmail is the noticiation medium of email.
	MediumEmail = Medium("email")
	// MediumMobilePush is the notification medium of a push notification.
	MediumMobilePush = Medium("mobile_push")
	// MediumSMS is the notification medium of SMS.
	MediumSMS = Medium("SMS")
)

// ListSettingsRequest is the request for ListSettings.
type ListSettingsRequest struct {
	ID           uint64
	UserID       uint64
	CanvasUserID uint64
	Type         Type
	Medium       Medium
}

// Setting represents a user's setting for one type of notifications on one medium.
type Setting struct {
	ID         uint64
	UserID     uint64
	Type       Type
	Medium     Medium
	Enabled    bool
	InsertedAt time.Time
}

// ListSettings lists notification settings.
func ListSettings(db services.DB, req *ListSettingsRequest) (*[]Setting, error) {
	q := util.Sq.
		Select(
			"notification_settings.id",
			"notification_settings.user_id",
			"notification_settings.notification_type_id",
			"notification_settings.medium",
			"notification_settings.enabled",
			"notification_settings.inserted_at",
		).
		From("notification_settings")

	if req.ID > 0 {
		q = q.Where(sq.Eq{"id": req.ID})
	}

	if req.UserID > 0 {
		q = q.Where(sq.Eq{"user_id": req.UserID})
	}

	if req.CanvasUserID > 0 {
		q = q.Join("users ON notification_settings.user_id = users.id").
			Where(sq.Eq{"users.canvas_user_id": req.CanvasUserID})
	}

	if req.Type != zeroType {
		q = q.Where(sq.Eq{"notification_type_id": req.Type})
	}

	if req.Medium != zeroMedium {
		q = q.Where(sq.Eq{"medium": req.Medium})
	}

	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building list notification settings sql: %w", err)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing list notification settings sql: %w", err)
	}

	defer rows.Close()

	var settings []Setting
	for rows.Next() {
		var s Setting
		err = rows.Scan(
			&s.ID,
			&s.UserID,
			&s.Type,
			&s.Medium,
			&s.Enabled,
			&s.InsertedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning list notification settings sql: %w", err)
		}

		settings = append(settings, s)
	}

	return &settings, nil
}
