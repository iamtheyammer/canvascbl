package db

import (
	"github.com/iamtheyammer/canvascbl/backend/src/db/canvas/courses"
	coursessvc "github.com/iamtheyammer/canvascbl/backend/src/db/services/courses"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
)

// UpsertMultipleCourses takes a course response from Canvas and upserts the courses
func UpsertMultipleCourses(cr *string) {
	ccr, err := courses.FromJSON(cr)
	if err != nil {
		handleError(err)
		return
	}

	var urs []coursessvc.UpsertRequest

	for _, c := range *ccr {
		urs = append(urs, coursessvc.UpsertRequest{
			Name:       c.Name,
			CourseCode: c.CourseCode,
			State:      c.State,
			UUID:       c.UUID,
			CourseID:   c.ID,
		})
	}

	err = coursessvc.UpsertMultiple(util.DB, &urs)
	if err != nil {
		handleError(errors.Wrap(err, "error upserting multiple courses"))
		return
	}

	return
}

// InsertMultipleAssignments takes an assignments response from Canvas and inserts the assignments
func InsertMultipleAssignments(ar *string) {
	car, err := courses.AssignmentsFromJSON(ar)
	if err != nil {
		handleError(errors.Wrap(err, "error getting CanvasAssignmentsResponse from JSON"))
		return
	}

	var air []coursessvc.AssignmentInsertRequest

	for _, a := range *car {
		air = append(air, coursessvc.AssignmentInsertRequest{
			CourseID: a.CourseID,
			CanvasID: a.ID,
			IsQuiz:   a.IsQuizAssignment,
			Name:     a.Name,
		})
	}

	err = coursessvc.InsertMultipleAssignments(util.DB, &air)
	if err != nil {
		handleError(errors.Wrap(err, "error inserting multiple assignments"))
		return
	}

	return
}
