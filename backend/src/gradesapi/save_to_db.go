package gradesapi

import (
	"fmt"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/courses"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/grades"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/outcomes"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/users"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"strconv"
	"strings"
)

func saveProfileToDB(p *canvasUserProfile) {
	err := users.UpsertProfile(db, &users.UpsertRequest{
		Name:         p.Name,
		Email:        p.PrimaryEmail,
		LTIUserID:    p.LtiUserID,
		CanvasUserID: int64(p.ID),
	})
	if err != nil {
		util.HandleError(fmt.Errorf("error saving user profile to db: %w", err))
		return
	}
}

func saveCoursesToDB(cs *[]canvasCourse) {
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

	err := courses.UpsertMultiple(db, &req)
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
						CanvasUserID: o.ID,
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
				CanvasUserID: o.ID,
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
		err := users.InsertUserObservees(trx, &users.UpsertObserveesRequest{
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

func saveOutcomeResultsToDB(results processedOutcomeResults) {
	var req []courses.OutcomeResultInsertRequest
	for cID, us := range results {
		for uID, os := range us {
			for oID, res := range os {
				for _, r := range res {
					aID, err := strconv.Atoi(strings.TrimPrefix(r.Links.Assignment, "assignment_"))
					if err != nil {
						util.HandleError(fmt.Errorf("failed to strip and convert a linked assignment id in an outcome result: %w", err))
						return
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

	err := courses.InsertMultipleOutcomeResults(db, &req)
	if err != nil {
		util.HandleError(fmt.Errorf("error inserting multiple outcome rollups: %w", err))
		return
	}
}

func saveGradesToDB(grds detailedGrades) {
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

	go func(request *[]grades.InsertRequest) {
		err := grades.Insert(db, request)
		if err != nil {
			util.HandleError(fmt.Errorf("error inserting grades: %w", err))
			return
		}
	}(&req)

	go func(request *[]courses.OutcomeRollupInsertRequest) {
		err := courses.InsertMultipleOutcomeRollups(db, request)
		if err != nil {
			util.HandleError(fmt.Errorf("error inserting multiple outcome averages (outcome rollups): %w", err))
			return
		}
	}(&rs)

	return
}

func saveAssignmentsToDB(ass []canvasAssignment, courseID string) {
	cID, err := strconv.Atoi(courseID)
	if err != nil {
		util.HandleError(fmt.Errorf("error converting course ID %s into an int: %w", courseID, err))
		return
	}

	var req []courses.AssignmentInsertRequest
	for _, a := range ass {
		req = append(req, courses.AssignmentInsertRequest{
			CourseID: uint64(cID),
			CanvasID: a.ID,
			IsQuiz:   a.IsQuizAssignment,
			Name:     a.Name,
		})
	}

	err = courses.InsertMultipleAssignments(db, &req)
	if err != nil {
		util.HandleError(fmt.Errorf("error inserting multiple assignments for course %d: %w", courseID, err))
		return
	}
}

func saveOutcomeToDB(o *canvasOutcome) {
	req := outcomes.InsertRequest{
		CanvasID:       o.ID,
		CourseID:       &o.ContextID,
		ContextID:      o.ContextID,
		DisplayName:    o.DisplayName,
		Title:          o.Title,
		MasteryPoints:  o.MasteryPoints,
		PointsPossible: o.PointsPossible,
	}

	if o.ContextType != "Course" {
		req.CourseID = nil
	}

	err := outcomes.UpsertOutcome(db, &req)
	if err != nil {
		util.HandleError(fmt.Errorf("error saving outcome %d to db: %w", o.ID, err))
		return
	}
}
