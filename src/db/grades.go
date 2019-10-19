package db

import (
	"github.com/iamtheyammer/canvascbl/backend/src/db/canvas/grades"
	gradessvc "github.com/iamtheyammer/canvascbl/backend/src/db/services/grades"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
	"strconv"
)

// InsertGrade inserts a grade from a outcome_rollups canvas API response
func InsertGrade(rr *string, courseID *string, userID *string) {
	cID, err := strconv.Atoi(*courseID)
	if err != nil {
		handleError(errors.Wrap(err, "error converting courseID to int"))
		return
	}

	uID, err := strconv.Atoi(*userID)
	if err != nil {
		handleError(errors.Wrap(err, "error converting userID to int"))
		return
	}

	crr, err := grades.GetCanvasRollupsResponseFromJsonString(rr)
	if err != nil {
		handleError(errors.Wrap(err, "error getting CanvasRollupsResponse from JSON string"))
		return
	}

	os, err := grades.GetOutcomeScoresFromCanvasRollupsResponse(crr)
	if err != nil {
		handleError(errors.Wrap(err, "error getting outcome scores from CanvasRollupsResponse"))
		return
	}

	// no graded outcomes for this class
	if len(*os) == 0 {
		return
	}

	grade := grades.CalculateGradeFromOutcomeScores(*os)

	err = gradessvc.Insert(util.DB, &gradessvc.InsertRequest{
		Grade:    grade,
		CourseID: cID,
		UserID:   uID,
	})

	if err != nil {
		handleError(errors.Wrap(err, "database error when inserting grades"))
		return
	}

	return
}
