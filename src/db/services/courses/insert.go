package courses

import (
	"database/sql"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
)

type AssignmentInsertRequest struct {
	CourseID uint64
	CanvasID uint64
	IsQuiz   bool
	Name     string
}

func InsertMultipleAssignments(db *sql.DB, req *[]AssignmentInsertRequest) error {
	q := util.Sq.
		Insert("assignments").
		Columns(
			"course_id",
			"canvas_id",
			"is_quiz",
			"name",
		).
		Suffix("ON CONFLICT DO NOTHING")

	for _, a := range *req {
		q = q.Values(
			a.CourseID,
			a.CanvasID,
			a.IsQuiz,
			a.Name,
		)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return errors.Wrap(err, "error building insert multiple assignments sql")
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "error executing insert multiple assignments sql")
	}

	return nil
}
