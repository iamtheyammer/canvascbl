package notifications

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"time"
)

// Medium represents a notification medium.
type Medium string

// IsValid determines if a Medium is a valid Medium
func (m Medium) IsValid() bool {
	switch m {
	case MediumEmail:
	default:
		return false
	}

	return true
}

const (
	// GradeChangeNotificationID is the ID for the grade_change notification type.
	TypeGradeChange = 1

	// zeroMedium is Medium's zero value
	zeroMedium = Medium("")
	// MediumEmail is the noticiation medium of email.
	MediumEmail = Medium("email")
	// MediumMobilePush is the notification medium of a push notification.
	MediumMobilePush = Medium("mobile_push")
	// MediumSMS is the notification medium of SMS.
	MediumSMS = Medium("SMS")
)

// Setting represents a user's setting for one type of notifications on one medium.
type Setting struct {
	ID           uint64
	UserID       uint64
	CanvasUserID uint64
	Type         uint64
	Medium       Medium
	InsertedAt   time.Time
}

// Type represents a notification type.
type Type struct {
	ID          uint64
	Name        string
	ShortName   string
	Description string
	InsertedAt  string
}

// ListTypesRequest is the request for ListTypes. It's empty, but it can be added to in the future.
type ListTypesRequest struct{}

// ListSettingsRequest is the request for ListSettings.
type ListSettingsRequest struct {
	ID           uint64
	UserID       uint64
	CanvasUserID uint64
	Type         uint64
	Medium       Medium
}

// ListTypes lists notification types.
func ListTypes(db services.DB, req *ListTypesRequest) (*[]Type, error) {
	query, args, err := util.Sq.
		Select(
			"id",
			"name",
			"short_name",
			"description",
			"inserted_at",
		).
		From("notification_types").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building list notification types sql: %w", err)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing list notification types sql: %w", err)
	}

	defer rows.Close()

	var ts []Type
	for rows.Next() {
		var (
			t           Type
			description sql.NullString
		)
		err = rows.Scan(
			&t.ID,
			&t.Name,
			&t.ShortName,
			&description,
			&t.InsertedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning list notification types sql: %w", err)
		}

		if description.Valid {
			t.Description = description.String
		}

		ts = append(ts, t)
	}

	return &ts, nil
}

// ListSettings lists notification settings.
func ListSettings(db services.DB, req *ListSettingsRequest) (*[]Setting, error) {
	q := util.Sq.
		Select(
			"notification_settings.id",
			"notification_settings.user_id",
			"users.canvas_user_id",
			"notification_settings.notification_type_id",
			"notification_settings.medium",
			"notification_settings.inserted_at",
		).
		From("notification_settings").
		Join("users ON notification_settings.user_id = users.id")

	if req.ID > 0 {
		q = q.Where(sq.Eq{"id": req.ID})
	}

	if req.UserID > 0 {
		q = q.Where(sq.Eq{"user_id": req.UserID})
	}

	if req.CanvasUserID > 0 {
		q = q.Where(sq.Eq{"users.canvas_user_id": req.CanvasUserID})
	}

	if req.Type > 0 {
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
			&s.CanvasUserID,
			&s.Type,
			&s.Medium,
			&s.InsertedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning list notification settings sql: %w", err)
		}

		settings = append(settings, s)
	}

	return &settings, nil
}
