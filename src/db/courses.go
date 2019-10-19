package db

import (
	"fmt"
	"github.com/iamtheyammer/canvascbl/backend/src/db/canvas/courses"
	coursessvc "github.com/iamtheyammer/canvascbl/backend/src/db/services/courses"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

// UpsertMultipleCourses takes a course response from Canvas and upserts the courses
func UpsertMultipleCourses(cr *string) {
	ccr, err := courses.FromJSON(cr)
	if err != nil {
		handleError(err)
		return
	}

	if len(*ccr) < 1 {
		// user has no courses
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

	if len(*car) < 1 {
		// no assignments
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

// InsertMultipleOutcomeRollups inserts multiple outcome rollups from string JSON
func InsertMultipleOutcomeRollups(orj *string, courseID *string) {
	corr, err := courses.OutcomeRollupsFromJSON(orj)
	if err != nil {
		handleError(errors.Wrap(err, "error getting CanvasOutcomeRollupsResponse from JSON"))
		return
	}

	if len(corr.Rollups) < 1 {
		handleError(errors.New("less than 1 rollup returned from OutcomeRollupsFromJSON"))
		return
	}

	// one per user, but we do this on a single user basis
	rollup := corr.Rollups[0]

	if len(rollup.Scores) < 1 {
		// no scores return
		return
	}

	cID, err := strconv.Atoi(*courseID)
	if err != nil {
		handleError(errors.Wrap(err, "error converting section id to int"))
		return
	}

	uID, err := strconv.Atoi(rollup.Links.User)
	if err != nil {
		handleError(errors.Wrap(err, "error converting user id to int"))
		return
	}

	var orir []coursessvc.OutcomeRollupInsertRequest

	for _, s := range rollup.Scores {
		oID, err := strconv.Atoi(s.Links.Outcome)
		if err != nil {
			handleError(errors.Wrap(err, fmt.Sprintf("error converting %s to an int", s.Links.Outcome)))
			return
		}

		orir = append(orir, coursessvc.OutcomeRollupInsertRequest{
			CourseID:      uint64(cID),
			OutcomeID:     uint64(oID),
			Score:         s.Score,
			TimesAssessed: uint64(s.Count),
		})
	}

	err = coursessvc.InsertMultipleOutcomeRollups(util.DB, uint64(uID), &orir)
	if err != nil {
		handleError(errors.Wrap(err, "error inserting multiple outcome rollups"))
		return
	}

	return
}

func InsertMultipleOutcomeResults(orj *string, courseID *string) {
	cID, err := strconv.Atoi(*courseID)
	if err != nil {
		handleError(errors.Wrap(err, "error converting courseID to int"))
		return
	}

	corr, err := courses.OutcomeResultsFromJSON(orj)
	if err != nil {
		handleError(errors.Wrap(err, "error unmarshaling into CanvasOutcomeResultsResponse"))
		return
	}

	var orir []coursessvc.OutcomeResultInsertRequest

	for _, r := range corr.OutcomeResults {
		aID, err := strconv.Atoi(strings.TrimPrefix(r.Links.Assignment, "assignment_"))
		if err != nil {
			handleError(errors.Wrap(err, "failed to strip and convert a linked assignment id in an outcome result"))
			return
		}

		oID, err := strconv.Atoi(r.Links.LearningOutcome)
		if err != nil {
			handleError(errors.Wrap(err, "failed to convert a linked learning outcome id in an outcome result"))
			return
		}

		uID, err := strconv.Atoi(r.Links.User)
		if err != nil {
			handleError(errors.Wrap(err, "failed to convert a linked user id in an outcome result"))
			return
		}

		orir = append(orir, coursessvc.OutcomeResultInsertRequest{
			ID:              r.ID,
			CourseID:        uint64(cID),
			AssignmentID:    uint64(aID),
			OutcomeID:       uint64(oID),
			UserID:          uint64(uID),
			AchievedMastery: r.Mastery,
			Score:           r.Score,
			Possible:        r.Possible,
			SubmissionTime:  r.SubmittedOrAssessedAt,
		})
	}

	err = coursessvc.InsertMultipleOutcomeResults(util.DB, &orir)
	if err != nil {
		handleError(errors.Wrap(err, "error inserting multiple outcome results"))
		return
	}

	return
}
