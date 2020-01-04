package courses

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
	"time"
)

type Course struct {
	ID         uint64
	Name       string
	CourseCode string
	State      string
	UUID       string
	CourseID   uint64
	InsertedAt time.Time
}

// GetForUser returns all courses the user has a grade in
func GetForUser(db services.DB, userIDs []uint64) (*[]Course, error) {
	query, args, err := util.Sq.
		Select(
			"courses.id",
			"courses.name",
			"courses.course_code",
			"courses.state",
			"courses.uuid",
			"courses.course_id",
			"courses.inserted_at",
		).
		From("courses").
		Distinct().
		LeftJoin("grades ON grades.course_id = courses.course_id").
		Where(sq.Eq{"grades.user_canvas_id": userIDs}).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "error building get courses for user sql")
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "error querying get courses for user")
	}

	defer rows.Close()

	var courses []Course

	for rows.Next() {
		var (
			c          Course
			courseCode sql.NullString
		)

		err := rows.Scan(
			&c.ID,
			&c.Name,
			&courseCode,
			&c.State,
			&c.UUID,
			&c.CourseID,
			&c.InsertedAt,
		)
		if err != nil {
			return nil, errors.Wrap(err, "error scanning get courses for user")
		}

		if courseCode.Valid {
			c.CourseCode = courseCode.String
		} else {
			c.CourseCode = c.Name
		}

		courses = append(courses, c)
	}

	return &courses, nil
}
