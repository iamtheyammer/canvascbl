package grades

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
	"time"
)

type ListRequest struct {
	UserCanvasIDs *[]uint64
	Before        *time.Time
	After         *time.Time
	CourseIDs     *[]uint64
	ManualFetch   *bool
}

// ListDistanceLearningRequest is a request for ListDistanceLearning.
type ListDistanceLearningRequest struct {
	DistanceLearningCourseID uint64
	OriginalCourseID         uint64
	UserCanvasIDs            []uint64
	Before                   time.Time
	After                    time.Time
}

type Grade struct {
	ID           uint64
	CourseID     uint64
	Grade        string
	UserCanvasID uint64
	ManualFetch  bool
	InsertedAt   time.Time
}

// DistanceLearningGrade represents a grade from distance_learning_grades.
type DistanceLearningGrade struct {
	ID                       uint64
	DistanceLearningCourseID uint64
	OriginalCourseID         uint64
	Grade                    string
	UserCanvasID             uint64
	ManualFetch              bool
	InsertedAt               time.Time
}

/*
List lists stored grades.

You will only get one row per course and user together.
*/
func List(db services.DB, req *ListRequest) (*[]Grade, error) {
	q := util.Sq.
		Select(
			"DISTINCT ON (grades.course_id, grades.user_canvas_id) id",
			"user_canvas_id",
			"course_id",
			"grade",
			"manual_fetch",
			"inserted_at",
		).
		From("grades").
		OrderBy("grades.course_id, grades.user_canvas_id, grades.inserted_at DESC").
		Where(sq.Eq{"user_canvas_id": req.UserCanvasIDs})

	if req.Before != nil {
		// using this weird workaround because it doesn't work any other way
		q = q.Where("inserted_at < to_timestamp(?)", req.Before.Unix())
	}

	if req.After != nil {
		q = q.Where(sq.Gt{"inserted_at": req.After})
	}

	if req.CourseIDs != nil {
		q = q.Where(sq.Eq{"course_id": req.CourseIDs})
	}

	if req.ManualFetch != nil {
		q = q.Where(sq.Eq{"manual_fetch": *req.ManualFetch})
	}

	query, args, err := q.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "error building list grades sql")
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "error opening rows in list grades sql")
	}

	defer rows.Close()

	var grades []Grade

	for rows.Next() {
		var g Grade

		err := rows.Scan(
			&g.ID,
			&g.UserCanvasID,
			&g.CourseID,
			&g.Grade,
			&g.ManualFetch,
			&g.InsertedAt,
		)

		if err != nil {
			return nil, errors.Wrap(err, "error scanning list grades sql")
		}

		grades = append(grades, g)
	}

	return &grades, nil
}

// ListDistanceLearning lists distance learning grades from distance_learning_grades.
// It returns the most recent grade per user.
func ListDistanceLearning(db services.DB, req *ListDistanceLearningRequest) (*[]DistanceLearningGrade, error) {
	q := util.Sq.
		Select(
			"DISTINCT ON (distance_learning_grades.original_course_id, "+
				"distance_learning_grades.user_canvas_id) id",
			"distance_learning_course_id",
			"original_course_id",
			"grade",
			"user_canvas_id",
			"manual_fetch",
			"inserted_at",
		).
		From("distance_learning_grades").
		OrderBy("distance_learning_grades.original_course_id, " +
			"distance_learning_grades.user_canvas_id, " +
			"distance_learning_grades.inserted_at DESC")

	if req.DistanceLearningCourseID > 0 {
		q = q.Where(sq.Eq{"distance_learning_course_id": req.DistanceLearningCourseID})
	}

	if req.OriginalCourseID > 0 {
		q = q.Where(sq.Eq{"original_course_id": req.OriginalCourseID})
	}

	if len(req.UserCanvasIDs) > 0 {
		q = q.Where(sq.Eq{"user_canvas_id": req.UserCanvasIDs})
	}

	if !req.Before.IsZero() {
		q = q.Where("inserted_at < to_timestamp(?)", req.Before.Unix())
	}

	if !req.After.IsZero() {
		q = q.Where("inserted_at > to_timestamp(?)", req.After.Unix())
	}

	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building list distance learning grades sql: %w", err)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing list distance learning grades sql: %w", err)
	}

	defer rows.Close()

	var dlGrades []DistanceLearningGrade
	for rows.Next() {
		var g DistanceLearningGrade
		err := rows.Scan(
			&g.ID,
			&g.DistanceLearningCourseID,
			&g.OriginalCourseID,
			&g.Grade,
			&g.UserCanvasID,
			&g.ManualFetch,
			&g.InsertedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning list distance learning grades sql: %w", err)
		}

		dlGrades = append(dlGrades, g)
	}

	return &dlGrades, nil
}
