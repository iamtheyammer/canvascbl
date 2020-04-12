package grades

import (
	"fmt"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
)

type InsertRequest struct {
	Grade        string
	ManualFetch  bool
	CourseID     int
	UserCanvasID int
}

type InsertDistanceLearningRequest struct {
	DistanceLearningCourseID uint64
	OriginalCourseID         uint64
	Grade                    string
	UserCanvasID             uint64
	ManualFetch              bool
}

func Insert(db services.DB, req *[]InsertRequest) error {
	q := util.Sq.
		Insert("grades").
		Columns("course_id", "grade", "user_canvas_id", "manual_fetch")

	for _, r := range *req {
		q = q.Values(
			r.CourseID,
			r.Grade,
			r.UserCanvasID,
			r.ManualFetch,
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

// InsertDistanceLearning inserts distance learning grades into distance_learning_grades
func InsertDistanceLearning(db services.DB, req *[]InsertDistanceLearningRequest) error {
	q := util.Sq.
		Insert("distance_learning_grades").
		Columns(
			"distance_learning_course_id",
			"original_course_id",
			"grade",
			"user_canvas_id",
			"manual_fetch",
		)

	for _, r := range *req {
		q = q.Values(
			r.DistanceLearningCourseID,
			r.OriginalCourseID,
			r.Grade,
			r.UserCanvasID,
			r.ManualFetch,
		)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return fmt.Errorf("error building insert distance learning grades sql: %w", err)
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("error executing insert distance learning grades sql: %w", err)
	}

	return nil
}
