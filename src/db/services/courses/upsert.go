package courses

import (
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
)

type UpsertRequest struct {
	Name       string
	CourseCode string
	State      string
	UUID       string
	CourseID   int64
}

// UpsertMultiple takes a course and if it already exists in the database, it ignores it (otherwise it's inserted)
func UpsertMultiple(db services.DB, c *[]UpsertRequest) error {
	q := util.Sq.
		Insert("courses").
		Columns(
			"name",
			"course_code",
			"state",
			"uuid",
			"course_id",
		).
		Suffix("ON CONFLICT DO NOTHING")

	for _, course := range *c {
		if course.Name != course.CourseCode {
			q = q.Values(
				course.Name,
				course.CourseCode,
				course.State,
				course.UUID,
				course.CourseID,
			)
			continue
		}

		q = q.Values(
			course.Name,
			nil,
			course.State,
			course.UUID,
			course.CourseID,
		)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return errors.Wrap(err, "error building upsert course sql")
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "error executing upsert course sql")
	}

	return nil
}
