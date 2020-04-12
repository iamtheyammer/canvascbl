package gradesapi

import (
	"fmt"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/courses"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/gpas"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/grades"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/outcomes"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/users"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"strconv"
	"strings"
)

func prepareProfileForDB(p *canvasUserProfile) *users.UpsertRequest {
	return &users.UpsertRequest{
		Name:         p.Name,
		Email:        p.PrimaryEmail,
		LTIUserID:    p.LtiUserID,
		CanvasUserID: int64(p.ID),
	}
}

func saveProfileToDB(p *canvasUserProfile) {
	_, err := users.UpsertProfile(db, prepareProfileForDB(p), false)
	if err != nil {
		util.HandleError(fmt.Errorf("error saving user profile to db: %w", err))
		return
	}
}

func prepareCoursesForDB(cs *[]canvasCourse) *[]courses.UpsertRequest {
	var req []courses.UpsertRequest
	for _, c := range *cs {
		req = append(req, courses.UpsertRequest{
			Name:       c.Name,
			CourseCode: c.CourseCode,
			State:      c.WorkflowState,
			UUID:       c.UUID,
			CourseID:   int64(c.ID),
		})
	}

	return &req
}

func saveCoursesToDB(cs *[]canvasCourse) {
	err := courses.UpsertMultiple(db, prepareCoursesForDB(cs))
	if err != nil {
		util.HandleError(fmt.Errorf("error saving courses to db: %w", err))
		return
	}
}

func saveObserveesToDB(cObs *[]canvasObservee, requestingUserID uint64) {
	var obs []users.Observee
	for _, o := range *cObs {
		obs = append(obs, users.Observee{
			ObserverUserID: requestingUserID,
			CanvasUserID:   o.ID,
			Name:           o.Name,
		})
	}

	// now, we'll start a db transaction
	trx, err := db.Begin()
	if err != nil {
		util.HandleError(fmt.Errorf("error beginning handle observees transaction: %w", err))
		return
	}

	// get the user's current observees
	dbObserveesP, err := users.ListObservees(trx, &users.ListObserveesRequest{ObserverCanvasUserID: requestingUserID})
	if err != nil {
		util.HandleError(fmt.Errorf("error listing user observees: %w", err))
		return
	}

	dbObservees := *dbObserveesP

	var (
		toSoftDelete, toUnSoftDelete []uint64
		toUpsert                     []users.Observee
	)

	for _, o := range obs {
		foundIDMatch := false
		for _, dbO := range dbObservees {
			if o.CanvasUserID == dbO.CanvasUserID {
				// if the names don't match, upsert
				if o.Name != dbO.Name {
					toUpsert = append(toUpsert, users.Observee{
						CanvasUserID: o.CanvasUserID,
						Name:         o.Name,
					})
				}

				// if it was previously deleted, undelete
				if !dbO.DeletedAt.IsZero() {
					toUnSoftDelete = append(toUnSoftDelete, dbO.CanvasUserID)
				}

				foundIDMatch = true
			}
		}

		// if it exists in observees from canvas but not in db, upsert.
		if !foundIDMatch {
			toUpsert = append(toUpsert, users.Observee{
				CanvasUserID: o.CanvasUserID,
				Name:         o.Name,
			})
		}
	}

	for _, dbO := range dbObservees {
		foundIDMatch := false
		for _, o := range obs {
			if dbO.CanvasUserID == o.CanvasUserID {
				foundIDMatch = true
			}
		}

		// if it exists in the db and it's not already deleted
		if !foundIDMatch && dbO.DeletedAt.IsZero() {
			toSoftDelete = append(toSoftDelete, dbO.CanvasUserID)
		}
	}

	if len(toSoftDelete) > 0 {
		err := users.SoftDeleteUserObservees(trx, toSoftDelete)
		if err != nil {
			util.HandleError(fmt.Errorf("error soft deleting user observees: %w", err))
			rollbackErr := trx.Rollback()
			if rollbackErr != nil {
				util.HandleError(fmt.Errorf("error rolling back handle observees trx: %w", rollbackErr))
			}
			return
		}
	}

	if len(toUnSoftDelete) > 0 {
		err := users.UnSoftDeleteUserObservees(trx, toUnSoftDelete)
		if err != nil {
			util.HandleError(fmt.Errorf("error un-soft deleting user observees: %w", err))
			rollbackErr := trx.Rollback()
			if rollbackErr != nil {
				util.HandleError(fmt.Errorf("error rolling back handle observees trx: %w", rollbackErr))
			}
			return
		}
	}

	if len(toUpsert) > 0 {
		err := users.UpsertUserObservees(trx, &users.UpsertObserveesRequest{
			Observees:            toUpsert,
			ObserverCanvasUserID: requestingUserID,
		})
		if err != nil {
			util.HandleError(fmt.Errorf("error upserting user observees: %w", err))
			rollbackErr := trx.Rollback()
			if rollbackErr != nil {
				util.HandleError(fmt.Errorf("error rolling back handle observees trx: %w", rollbackErr))
			}
			return
		}
	}

	err = trx.Commit()
	if err != nil {
		util.HandleError(fmt.Errorf("error committing handle observees trx: %w", err))
		rollbackErr := trx.Rollback()
		if rollbackErr != nil {
			util.HandleError(fmt.Errorf("error rolling back handle observees trx: %w", rollbackErr))
		}
		return
	}

	return
}

func prepareOutcomeResultsForDB(results processedOutcomeResults) (*[]courses.OutcomeResultInsertRequest, error) {
	var req []courses.OutcomeResultInsertRequest
	for cID, us := range results {
		for uID, os := range us {
			for oID, res := range os {
				for _, r := range res {
					aID, err := strconv.Atoi(strings.TrimPrefix(r.Links.Assignment, "assignment_"))
					if err != nil {
						return nil, fmt.Errorf("failed to strip and convert a linked assignment id in an outcome result: %w", err)
					}

					req = append(req, courses.OutcomeResultInsertRequest{
						ID:              r.ID,
						CourseID:        cID,
						AssignmentID:    uint64(aID),
						OutcomeID:       oID,
						UserID:          uID,
						AchievedMastery: r.Mastery,
						Score:           r.Score,
						Possible:        r.Possible,
						SubmissionTime:  r.SubmittedOrAssessedAt,
					})
				}
			}
		}
	}

	return &req, nil
}

func saveOutcomeResultsToDB(results processedOutcomeResults) {
	req, err := prepareOutcomeResultsForDB(results)
	if err != nil {
		util.HandleError(fmt.Errorf("error preparing outcome results for db: %w", err))
		return
	}

	err = courses.InsertMultipleOutcomeResults(db, req)
	if err != nil {
		util.HandleError(fmt.Errorf("error inserting multiple outcome rollups: %w", err))
		return
	}
}

func prepareGradesForDB(
	grds detailedGrades,
	manualFetch bool) (
	*[]grades.InsertRequest,
	*[]courses.OutcomeRollupInsertRequest,
) {
	var req []grades.InsertRequest
	var rs []courses.OutcomeRollupInsertRequest
	for uID, cs := range grds {
		for cID, grd := range cs {
			if grd.Grade == naGrade {
				continue
			}

			req = append(req, grades.InsertRequest{
				Grade:        grd.Grade.Grade,
				CourseID:     int(cID),
				UserCanvasID: int(uID),
				ManualFetch:  manualFetch,
			})

			for oID, avg := range grd.Averages {
				rs = append(rs, courses.OutcomeRollupInsertRequest{
					CanvasUserID: uID,
					CourseID:     cID,
					OutcomeID:    oID,
					Score:        avg.Average,
				})
			}
		}
	}

	return &req, &rs
}

func saveGradesToDB(grds detailedGrades, manualFetch bool) {
	req, rs := prepareGradesForDB(grds, manualFetch)

	go func(request *[]grades.InsertRequest) {
		err := grades.Insert(db, request)
		if err != nil {
			util.HandleError(fmt.Errorf("error inserting grades: %w", err))
			return
		}
	}(req)

	go func(request *[]courses.OutcomeRollupInsertRequest) {
		err := courses.InsertMultipleOutcomeRollups(db, request)
		if err != nil {
			util.HandleError(fmt.Errorf("error inserting multiple outcome averages (outcome rollups): %w", err))
			return
		}
	}(rs)

	return
}

func prepareAssignmentsForDB(ass []canvasAssignment, courseID string) (*[]courses.AssignmentUpsertRequest, error) {
	cID, err := strconv.Atoi(courseID)
	if err != nil {
		return nil, fmt.Errorf("error converting course ID %s into an int: %w", courseID, err)
	}

	var req []courses.AssignmentUpsertRequest
	for _, a := range ass {
		req = append(req, courses.AssignmentUpsertRequest{
			CourseID: uint64(cID),
			CanvasID: a.ID,
			IsQuiz:   a.IsQuizAssignment,
			Name:     a.Name,
			DueAt:    a.DueAt,
		})
	}

	return &req, nil
}

func saveAssignmentsToDB(ass []canvasAssignment, courseID string) {
	req, err := prepareAssignmentsForDB(ass, courseID)
	if err != nil {
		util.HandleError(fmt.Errorf("error preparing assignments for db: %w", err))
		return
	}

	err = courses.UpsertMultipleAssignments(db, req)
	if err != nil {
		util.HandleError(fmt.Errorf("error inserting multiple assignments for course %s: %w", courseID, err))
		return
	}
}

func prepareOutcomeForDB(o *canvasOutcome) *outcomes.InsertRequest {
	return &outcomes.InsertRequest{
		CanvasID:       o.ID,
		CourseID:       &o.ContextID,
		ContextID:      o.ContextID,
		DisplayName:    o.DisplayName,
		Title:          o.Title,
		MasteryPoints:  o.MasteryPoints,
		PointsPossible: o.PointsPossible,
	}
}

func saveOutcomeToDB(o *canvasOutcome) {
	req := prepareOutcomeForDB(o)

	if o.ContextType != "Course" {
		req.CourseID = nil
	}

	err := outcomes.UpsertOutcome(db, req)
	if err != nil {
		util.HandleError(fmt.Errorf("error saving outcome %d to db: %w", o.ID, err))
		return
	}
}

func prepareGPAForDB(g gpa, manualFetch bool) *[]gpas.InsertRequest {
	var req []gpas.InsertRequest

	for cuID, cGPA := range g {
		req = append(req, gpas.InsertRequest{
			CanvasUserID:     cuID,
			Weighted:         false,
			GPA:              cGPA.Unweighted.Default,
			GPAWithSubgrades: cGPA.Unweighted.Subgrades,
			ManualFetch:      manualFetch,
		})
	}

	return &req
}

func saveGPAToDB(g gpa, manualFetch bool) {
	err := gpas.InsertMultiple(db, prepareGPAForDB(g, manualFetch))
	if err != nil {
		util.HandleError(fmt.Errorf("error saving gpa to db: %w", err))
		return
	}
}

func prepareDistanceLearningGradesForDB(dlg distanceLearningGrades, manualFetch bool) *[]grades.InsertDistanceLearningRequest {
	var req []grades.InsertDistanceLearningRequest

	for uID, gs := range dlg {
		for _, g := range gs {
			req = append(req, grades.InsertDistanceLearningRequest{
				DistanceLearningCourseID: g.DistanceLearningCourseID,
				OriginalCourseID:         g.OriginalCourseID,
				Grade:                    g.Grade.Grade,
				UserCanvasID:             uID,
				ManualFetch:              manualFetch,
			})
		}
	}

	return &req
}

func saveDistanceLearningGradesToDB(dlg distanceLearningGrades, manualFetch bool) {
	err := grades.InsertDistanceLearning(db, prepareDistanceLearningGradesForDB(dlg, manualFetch))
	if err != nil {
		util.HandleError(fmt.Errorf("error saving distance learning grades to db: %w", err))
		return
	}
}
