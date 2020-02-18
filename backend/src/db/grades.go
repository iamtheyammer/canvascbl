package db

import (
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/courses"
	gradessvc "github.com/iamtheyammer/canvascbl/backend/src/db/services/grades"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
	"time"
)

var memoizedGradeAverages = map[uint64]struct {
	Result     float64
	NumInputs  uint64
	ValidUntil time.Time
}{}

func GetAverageGradeForCourse(courseID uint64) (*gradessvc.CourseGradeAverage, error) {
	return gradessvc.GetAverageForCourse(util.DB, courseID)
}

func GetMemoizedAverageGradeForCourse(courseID uint64, userIDs []uint64) (*float64, *uint64, error) {
	cs, err := courses.GetForUser(util.DB, userIDs)
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

func GetGradesForUserBeforeDate(userIDs []uint64, before time.Time) (*[]gradessvc.Grade, error) {
	mf := true
	gs, err := gradessvc.List(util.DB, &gradessvc.ListRequest{
		UserCanvasIDs: &userIDs,
		Before:        &before,
		// since this is used for previous grades, we only want grades the user fetched manually
		ManualFetch: &mf,
	})
	if err != nil {
		return nil, errors.Wrap(err, "error getting grades")
	}

	return gs, nil
}
