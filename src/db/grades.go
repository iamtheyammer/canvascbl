package db

import (
	"github.com/iamtheyammer/canvascbl/backend/src/db/canvas/grades"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/courses"
	gradessvc "github.com/iamtheyammer/canvascbl/backend/src/db/services/grades"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

var memoizedGradeAverages = map[uint64]struct {
	Result     float64
	NumInputs  uint64
	ValidUntil time.Time
}{}

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

	if len(grade) < 1 {
		return
	}

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

func GetAverageGradeForCourse(courseID uint64) (*gradessvc.CourseGradeAverage, error) {
	return gradessvc.GetAverageForCourse(util.DB, courseID)
}

func GetMemoizedAverageGradeForCourse(courseID uint64, userID uint64) (*float64, *uint64, error) {
	cs, err := courses.GetForUser(util.DB, userID)
	if err != nil {
		return nil, nil, nil
	}

	userHasCourse := false
	for _, c := range *cs {
		if c.CourseID == courseID {
			userHasCourse = true
			break
		}
	}

	if !userHasCourse {
		return nil, nil, nil
	}

	v, ok := memoizedGradeAverages[courseID]

	if !ok {
		avg, err := gradessvc.GetAverageForCourse(util.DB, courseID)
		if err != nil {
			return nil, nil, errors.Wrap(err, "error getting grade average for course")
		}
		memoizedGradeAverages[courseID] = struct {
			Result     float64
			NumInputs  uint64
			ValidUntil time.Time
		}{Result: avg.Average, ValidUntil: time.Now().Add(time.Minute * 5), NumInputs: avg.NumInputs}

		return &avg.Average, &avg.NumInputs, nil
	}

	if v.ValidUntil.Before(time.Now()) {
		avg, err := gradessvc.GetAverageForCourse(util.DB, courseID)
		if err != nil {
			return nil, nil, errors.Wrap(err, "error getting grade average for course")
		}
		memoizedGradeAverages[courseID] = struct {
			Result     float64
			NumInputs  uint64
			ValidUntil time.Time
		}{Result: avg.Average, ValidUntil: time.Now().Add(time.Minute * 5), NumInputs: avg.NumInputs}
	}

	v, _ = memoizedGradeAverages[courseID]

	return &v.Result, &v.NumInputs, nil
}
