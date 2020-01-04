package grades

import (
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
)

type InsertRequest struct {
	Grade        string
	CourseID     int
	UserCanvasID int
}

func Insert(db services.DB, req *[]InsertRequest) error {
	q := util.Sq.
		Insert("grades").
		Columns("course_id", "grade", "user_canvas_id")

	for _, r := range *req {
		q = q.Values(
			r.CourseID,
			r.Grade,
			r.UserCanvasID,
		)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return errors.Wrap(err, "error building insert grade sql")
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "error executing insert grade sql")
	}

	return nil
}
