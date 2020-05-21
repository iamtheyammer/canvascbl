package gradesapi

import (
	"fmt"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/canvas_tokens"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/courses"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/enrollments"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/gpas"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/grades"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/outcomes"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/submissions"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/users"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"strconv"
	"strings"
	"time"
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

func prepareCoursesForDB(cs *[]canvasCourse) (*[]courses.UpsertRequest, *[]enrollments.UpsertRequest) {
	var (
		cReq []courses.UpsertRequest
		eReq []enrollments.UpsertRequest
	)
	for _, c := range *cs {
		cReq = append(cReq, courses.UpsertRequest{
			Name:       c.Name,
			CourseCode: c.CourseCode,
			State:      c.WorkflowState,
			UUID:       c.UUID,
			CourseID:   int64(c.ID),
		})

		for _, e := range c.Enrollments {
			eReq = append(eReq, enrollments.UpsertRequest{
				CourseID:               c.ID,
				UserCanvasID:           e.UserID,
				AssociatedUserCanvasID: e.AssociatedUserID,
				Type:                   enrollments.Type(e.Type),
				Role:                   enrollments.Role(e.Role),
				State:                  e.EnrollmentState,
			})
		}
	}

	return &cReq, &eReq
}

func saveCoursesToDB(cs *[]canvasCourse) {
	cReqs, eReqs := prepareCoursesForDB(cs)

	trx, err := db.Begin()
	if err != nil {
		util.HandleError(fmt.Errorf("error beginning save courses to db transaction: %w", err))
		return
	}

	err = courses.UpsertMultiple(trx, cReqs)
	if err != nil {
		util.HandleError(fmt.Errorf("error saving courses to db: %w", err))

		rbErr := trx.Rollback()
		if rbErr != nil {
			util.HandleError(fmt.Errorf("error rolling back save courses to db transaction at courses: %w", err))
		}
		return
	}

	err = enrollments.Upsert(trx, eReqs)
	if err != nil {
		util.HandleError(fmt.Errorf("error saving enrollments (in courses) to db: %w", err))

		rbErr := trx.Rollback()
		if rbErr != nil {
			util.HandleError(fmt.Errorf("error rolling back save ️courses to db transaction at enrollments: %w"))
		}

		return
	}

	err = trx.Commit()
	if err != nil {
		util.HandleError(fmt.Errorf("error committing save courses transaction to db: %w", err))

		rbErr := trx.Rollback()
		if rbErr != nil {
			util.HandleError(fmt.Errorf("error rolling back save courses to db transaction at commit: %w", err))
		}

		return
	}

	return
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

	numReqs := len(*req)

	if numReqs < 1 {
		return
	}

	// chunking required? if NO:
	if numReqs < courses.MultipleOutcomeResultsUpsertChunkSize {
		err = courses.InsertMultipleOutcomeResults(db, req)
		if err != nil {
			util.HandleError(fmt.Errorf("error inserting multiple outcome rollups: %w", err))
			return
		}
		return
	}

	// chunking
	chunked := [][]courses.OutcomeResultInsertRequest{{}}
	curChunk := 0
	chunkLen := 0

	for _, r := range *req {
		// if we are over the max per chunk, move to next chunk
		if chunkLen > courses.MultipleOutcomeResultsUpsertChunkSize {
			curChunk++
			chunked = append(chunked, []courses.OutcomeResultInsertRequest{})
			chunkLen = 0
		}

		// add to the current chunk
		chunked[curChunk] = append(chunked[curChunk], r)

		// add to the number in the current chunk
		chunkLen++
	}

	trx, err := db.Begin()
	if err != nil {
		util.HandleError(fmt.Errorf("error beginning insert multiple outcome results trx: %w", err))
		return
	}

	for i, ch := range chunked {
		err := courses.InsertMultipleOutcomeResults(trx, &ch)
		if err != nil {
			util.HandleError(fmt.Errorf("error inserting multiple outcome results (chunk %d): %w", i, err))

			rbErr := trx.Rollback()
			if rbErr != nil {
				util.HandleError(fmt.Errorf("error rolling back insert multiple outcome results trx (chunk %d): %w", i, err))
			}
			return
		}
	}

	err = trx.Commit()
	if err != nil {
		util.HandleError(fmt.Errorf("error committing multiple outcome results trx: %w", err))

		rbErr := trx.Rollback()
		if rbErr != nil {
			util.HandleError(fmt.Errorf("error rolling back insert multiple outcome results trx at commit: %w", err))
		}
		return
	}

	return
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

	if len(*req) > 0 {
		go func(request *[]grades.InsertRequest) {
			err := grades.Insert(db, request)
			if err != nil {
				util.HandleError(fmt.Errorf("error inserting grades: %w", err))
				return
			}
		}(req)
	}

	if len(*rs) > 0 {
		go func(request *[]courses.OutcomeRollupInsertRequest) {
			err := courses.InsertMultipleOutcomeRollups(db, request)
			if err != nil {
				util.HandleError(fmt.Errorf("error inserting multiple outcome averages (outcome rollups): %w", err))
				return
			}
		}(rs)
	}

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
			if g.Grade.Grade == "N/A" {
				continue
			}

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
	req := prepareDistanceLearningGradesForDB(dlg, manualFetch)
	if len(*req) < 1 {
		return
	}

	err := grades.InsertDistanceLearning(db, req)
	if err != nil {
		util.HandleError(fmt.Errorf("error saving distance learning grades to db: %w", err))
		return
	}
}

func prepareCanvasOAuth2GrantForDB(grant *canvasTokenGrantResponse) *canvas_tokens.InsertRequest {
	expAt := time.Now().Add(time.Duration(grant.ExpiresIn) * time.Second)

	return &canvas_tokens.InsertRequest{
		CanvasUserID: grant.User.ID,
		Token:        grant.AccessToken,
		RefreshToken: grant.RefreshToken,
		ExpiresAt:    &expAt,
	}
}

func saveCanvasOAuth2GrantToDB(grant *canvasTokenGrantResponse) {
	err := canvas_tokens.Insert(db, prepareCanvasOAuth2GrantForDB(grant))
	if err != nil {
		util.HandleError(fmt.Errorf("error saving canvas oauth2 grant to db: %w", err))
		return
	}
}

func prepareEnrollmentsForDB(req []canvasFullEnrollment) *[]enrollments.UpsertRequest {
	var reqs []enrollments.UpsertRequest

	for _, e := range req {
		reqs = append(reqs, enrollments.UpsertRequest{
			CanvasID:               e.ID,
			CourseID:               e.CourseID,
			UserCanvasID:           e.UserID,
			AssociatedUserCanvasID: e.AssociatedUserID,
			Type:                   enrollments.Role(e.Type).ToType(),
			Role:                   enrollments.Role(e.Role),
			State:                  e.EnrollmentState,
			CreatedAt:              e.CreatedAt,
			UpdatedAt:              e.UpdatedAt,
		})
	}

	return &reqs
}

func saveEnrollmentsToDB(req []canvasFullEnrollment) {
	err := enrollments.Upsert(db, prepareEnrollmentsForDB(req))
	if err != nil {
		util.HandleError(fmt.Errorf("error saving enrollments to db: %w", err))
		return
	}
}

func prepareSubmissionsForDB(req []canvasSubmission, courseID uint64) (*[]submissions.UpsertRequest, *[]submissions.AttachmentUpsertRequest) {
	var ss []submissions.UpsertRequest
	var as []submissions.AttachmentUpsertRequest
	for _, s := range req {
		ss = append(ss, submissions.UpsertRequest{
			CanvasID:         s.ID,
			CourseID:         courseID,
			AssignmentID:     s.AssignmentID,
			UserCanvasID:     s.UserID,
			Attempt:          s.Attempt,
			Score:            s.Score,
			WorkflowState:    submissions.WorkflowState(s.WorkflowState),
			GraderID:         s.GraderID,
			GradedAt:         s.GradedAt,
			Type:             s.SubmissionType,
			SubmittedAt:      s.SubmittedAt,
			HTMLURL:          s.URL,
			Late:             s.Late,
			Excused:          s.Excused,
			Missing:          s.Missing,
			LatePolicyStatus: s.LatePolicyStatus,
			PointsDeducted:   s.PointsDeducted,
			SecondsLate:      s.SecondsLate,
			ExtraAttempts:    s.ExtraAttempts,
			PostedAt:         s.PostedAt,
		})

		for _, a := range s.Attachments {
			as = append(as, submissions.AttachmentUpsertRequest{
				CanvasID:     a.ID,
				SubmissionID: s.ID,
				DisplayName:  a.DisplayName,
				Filename:     a.Filename,
				ContentType:  a.ContentType,
				URL:          a.URL,
				Size:         a.Size,
				CreatedAt:    &a.CreatedAt,
			})
		}
	}

	return &ss, &as
}

func saveSubmissionsToDB(req []canvasSubmission, courseID uint64) {
	ss, as := prepareSubmissionsForDB(req, courseID)

	go func() {
		err := submissions.Upsert(db, ss)
		if err != nil {
			util.HandleError(fmt.Errorf("error inserting submissions to db: %w", err))
		}
	}()

	if len(*as) > 0 {
		go func() {
			err := submissions.UpsertAttachments(db, as)
			if err != nil {
				util.HandleError(fmt.Errorf("error inserting submission attachments to db: %w", err))
			}
		}()
	}

}

// handleBatchGradesDBRequests handles tons of UserGradesDBRequests from fetch_all. It handles chunking and more.
func handleBatchGradesDBRequests(dbReqs []UserGradesDBRequests) {
	go func() {
		var (
			profiles []users.UpsertRequest
			// keeps out duplicates
			coursesMap = make(map[int64]struct{})
			cs         []courses.UpsertRequest
			// keeps out duplicates
			outcomeResultsMap = make(map[uint64]struct{})
			// chunked in 7281 due to postgres's 65535 parameter limit
			// 9 params each
			chunkedOutcomeResults           = [][]courses.OutcomeResultInsertRequest{{}}
			currentOutcomeResultChunk       = 0
			currentOutcomeResultChunkLength = 0
			// chunk in 16383 due to postgres's 65535 parameter limit
			// 4 params each
			chunkedGrades            = [][]grades.InsertRequest{{}}
			currentGradesChunk       = 0
			currentGradesChunkLength = 0
			// chunk in 13107 due to postgres's 65535 parameter limit
			// 5 params each
			chunkedRollupScores           = [][]courses.OutcomeRollupInsertRequest{{}}
			currentRollupScoreChunk       = 0
			currentRollupScoreChunkLength = 0

			gpaReqs                    []gpas.InsertRequest
			distanceLearningGradesReqs []grades.InsertDistanceLearningRequest

			// chunked in 7281 due to postgres's 65535 parameter limit
			// 9 params each
			chunkedEnrollments            = [][]enrollments.UpsertRequest{{}}
			currentEnrollmentsChunk       = 0
			currentEnrollmentsChunkLength = 0
			// enrollmentsCache makes sure we don't get the "can't DO UPDATE more than once" error
			enrollmentsCache = make(map[uint64]map[uint64]struct{})

			chunkedAssignments            = [][]courses.AssignmentUpsertRequest{{}}
			currentAssignmentsChunk       = 0
			currentAssignmentsChunkLength = 0
			assignmentsCache              = make(map[uint64]struct{})

			chunkedSubmissions            = [][]submissions.UpsertRequest{{}}
			currentSubmissionsChunk       = 0
			currentSubmissionsChunkLength = 0
			submissionsCache              = map[uint64]struct{}{}

			chunkedAttachments            = [][]submissions.AttachmentUpsertRequest{{}}
			currentAttachmentsChunk       = 0
			currentAttachmentsChunkLength = 0
			attachmentsCache              = map[uint64]struct{}{}
		)

		for _, r := range dbReqs {
			if r.Profile != nil {
				profiles = append(profiles, *r.Profile)
			}

			if r.Courses != nil {
				// keeps out duplicate courses
				for _, c := range *r.Courses {
					if _, ok := coursesMap[c.CourseID]; !ok {
						cs = append(cs, c)
						coursesMap[c.CourseID] = struct{}{}
					}
				}
			}

			if r.OutcomeResults != nil {
				for _, or := range *r.OutcomeResults {
					// keeps out duplicate outcome results
					if _, ok := outcomeResultsMap[or.ID]; !ok {
						// if we are over the max per chunk, move to next chunk
						if currentOutcomeResultChunkLength >= 7281 {
							currentOutcomeResultChunk++
							chunkedOutcomeResults = append(chunkedOutcomeResults, []courses.OutcomeResultInsertRequest{})
							currentOutcomeResultChunkLength = 0
						}

						// add to the current chunk
						chunkedOutcomeResults[currentOutcomeResultChunk] =
							append(chunkedOutcomeResults[currentOutcomeResultChunk], or)

						// no duplicates
						outcomeResultsMap[or.ID] = struct{}{}

						// add to the number in the current chunk
						currentOutcomeResultChunkLength++
					}
				}
			}

			if r.Grades != nil {
				for _, g := range *r.Grades {
					// if we are over the max per chunk, move to next chunk
					if currentGradesChunkLength >= 16383 {
						currentGradesChunk++
						chunkedGrades = append(chunkedGrades, []grades.InsertRequest{})
						currentGradesChunkLength = 0
					}

					// add to the current chunk
					chunkedGrades[currentGradesChunk] =
						append(chunkedGrades[currentGradesChunk], g)

					// add to the number in the current chunk
					currentGradesChunkLength++
				}
			}

			if r.RollupScores != nil {
				for _, rs := range *r.RollupScores {
					// if we are over the max per chunk, move to next chunk
					if currentRollupScoreChunkLength >= 13107 {
						currentRollupScoreChunk++
						chunkedRollupScores = append(chunkedRollupScores, []courses.OutcomeRollupInsertRequest{})
						currentRollupScoreChunkLength = 0
					}

					// add to the current chunk
					chunkedRollupScores[currentRollupScoreChunk] =
						append(chunkedRollupScores[currentRollupScoreChunk], rs)

					// add to the number in the current chunk
					currentRollupScoreChunkLength++
				}
			}

			if r.GPA != nil {
				gpaReqs = append(gpaReqs, *r.GPA...)
			}

			if r.DistanceLearningGrades != nil {
				distanceLearningGradesReqs = append(distanceLearningGradesReqs, *r.DistanceLearningGrades...)
			}

			if r.Enrollments != nil {
				for _, es := range r.Enrollments {
					if _, ok := enrollmentsCache[es.UserCanvasID]; ok {
						// user => course
						if enrollmentsCache[es.UserCanvasID] == nil {
							enrollmentsCache[es.UserCanvasID] = map[uint64]struct{}{}
						}

						if _, ok := enrollmentsCache[es.UserCanvasID][es.CourseID]; ok {
							continue
						} else {
							enrollmentsCache[es.UserCanvasID][es.CourseID] = struct{}{}
						}
					} else {
						enrollmentsCache[es.UserCanvasID] = map[uint64]struct{}{es.CourseID: {}}
					}

					// if we are over the max per chunk, move to next chunk
					if currentEnrollmentsChunkLength >= 7281 {
						currentEnrollmentsChunk++
						chunkedEnrollments = append(chunkedEnrollments, []enrollments.UpsertRequest{})
						currentEnrollmentsChunkLength = 0
					}

					// add to the current chunk
					chunkedEnrollments[currentEnrollmentsChunk] =
						append(chunkedEnrollments[currentEnrollmentsChunk], es)

					// add to the number in the current chunk
					currentEnrollmentsChunkLength++
				}
			}

			if r.Assignments != nil {
				for _, a := range r.Assignments {
					if _, ok := assignmentsCache[a.CanvasID]; ok {
						continue
					} else {
						assignmentsCache[a.CanvasID] = struct{}{}
					}

					// if we are over the max per chunk, move to next chunk
					if currentAssignmentsChunkLength >= courses.MultipleAssignmentsChunkSize {
						currentAssignmentsChunk++
						chunkedAssignments = append(chunkedAssignments, []courses.AssignmentUpsertRequest{})
						currentAssignmentsChunkLength = 0
					}

					// add to the current chunk
					chunkedAssignments[currentAssignmentsChunk] =
						append(chunkedAssignments[currentAssignmentsChunk], a)

					// add to the number in the current chunk
					currentAssignmentsChunkLength++
				}
			}

			if r.Submissions != nil {
				for _, s := range r.Submissions {
					if _, ok := submissionsCache[s.CanvasID]; ok {
						continue
					} else {
						submissionsCache[s.CanvasID] = struct{}{}
					}

					// if we are over the max per chunk, move to next chunk
					if currentSubmissionsChunkLength >= submissions.UpsertChunkSize {
						currentSubmissionsChunk++
						chunkedSubmissions = append(chunkedSubmissions, []submissions.UpsertRequest{})
						currentSubmissionsChunkLength = 0
					}

					// add to the current chunk
					chunkedSubmissions[currentSubmissionsChunk] =
						append(chunkedSubmissions[currentSubmissionsChunk], s)

					// add to the number in the current chunk
					currentSubmissionsChunkLength++
				}
			}

			if r.SubmissionAttachments != nil {
				for _, a := range r.SubmissionAttachments {
					if _, ok := attachmentsCache[a.CanvasID]; ok {
						continue
					} else {
						attachmentsCache[a.CanvasID] = struct{}{}
					}

					// if we are over the max per chunk, move to next chunk
					if currentAttachmentsChunk >= submissions.AttachmentsUpsertChunkSize {
						currentAttachmentsChunk++
						chunkedAttachments = append(chunkedAttachments, []submissions.AttachmentUpsertRequest{})
						currentAttachmentsChunkLength = 0
					}

					// add to the current chunk
					chunkedAttachments[currentAttachmentsChunk] =
						append(chunkedAttachments[currentAttachmentsChunk], a)

					// add to the number in the current chunk
					currentAttachmentsChunkLength++
				}
			}
		}

		// we'll request one at a time-- these are big, big requests
		trx, err := db.Begin()
		if err != nil {
			util.HandleError(fmt.Errorf("error beginning insert grades fetch_all data transaction: %w", err))
			return
		}

		rb := func(at string) {
			err := trx.Rollback()
			if err != nil {
				util.HandleError(
					fmt.Errorf("error rolling back insert grades fetch_all data transaction at %s: %w", at, err),
				)
			}
		}

		_, err = users.UpsertMultipleProfiles(trx, &profiles, false)
		if err != nil {
			rb("profiles")
			util.HandleError(fmt.Errorf("error upserting multiple profiles in insert fetch_all data: %w", err))
			return
		}

		err = courses.UpsertMultiple(trx, &cs)
		if err != nil {
			rb("courses")
			util.HandleError(fmt.Errorf("error upserting multiple courses in insert fetch_all data: %w", err))
			return
		}

		for i, req := range chunkedOutcomeResults {
			if len(req) < 1 {
				continue
			}

			err = courses.InsertMultipleOutcomeResults(trx, &req)
			if err != nil {
				rb("outcome results")
				util.HandleError(
					fmt.Errorf("error inserting multiple outcome results in insert fetch_all data (chunk %d): %w", i, err),
				)
				return
			}
		}

		for i, req := range chunkedGrades {
			if len(req) < 1 {
				continue
			}

			err = grades.Insert(trx, &req)
			if err != nil {
				rb("grades")
				util.HandleError(
					fmt.Errorf("error inserting multiple grades in insert fetch_all data (chunk %d): %w", i, err),
				)
				return
			}
		}

		for i, req := range chunkedRollupScores {
			if len(req) < 1 {
				continue
			}

			err = courses.InsertMultipleOutcomeRollups(trx, &req)
			if err != nil {
				rb("rollup scores")
				util.HandleError(
					fmt.Errorf("error inserting multiple rollup scores in insert fetch_all data (chunk %d): %w", i, err),
				)
				return
			}
		}

		//err = gpas.InsertMultiple(trx, &gpaReqs)
		//if err != nil {
		//	rb("gpas")
		//	util.HandleError(fmt.Errorf("error inserting multiple gpas in insert fetch_all data: %w", err))
		//	return
		//}

		if len(distanceLearningGradesReqs) > 0 {
			err = grades.InsertDistanceLearning(trx, &distanceLearningGradesReqs)
			if err != nil {
				rb("distance learning grades")
				util.HandleError(
					fmt.Errorf("error inserting distance learning grades in insert fetch_all data: %w", err),
				)
				return
			}
		}

		for i, req := range chunkedEnrollments {
			if len(req) < 1 {
				continue
			}

			err = enrollments.Upsert(trx, &req)
			if err != nil {
				rb("enrollments")
				util.HandleError(
					fmt.Errorf("error inserting multiple enrollments in insert fetch_all data (chunk %d): %w", i, err),
				)
				return
			}
		}

		for i, req := range chunkedAssignments {
			if len(req) < 1 {
				continue
			}

			err = courses.UpsertMultipleAssignments(trx, &req)
			if err != nil {
				rb("assignments")
				util.HandleError(
					fmt.Errorf("error inserting multiple assignments in insert fetch_all data (chunk %d): %w", i, err),
				)
				return
			}
		}

		for i, req := range chunkedSubmissions {
			if len(req) < 1 {
				continue
			}

			err = submissions.Upsert(trx, &req)
			if err != nil {
				rb("submissions")
				util.HandleError(
					fmt.Errorf("error inserting multiple submissions in insert fetch_all data (chunk %d): %w", i, err),
				)
				return
			}
		}

		for i, req := range chunkedAttachments {
			if len(req) < 1 {
				continue
			}

			err = submissions.UpsertAttachments(trx, &req)
			if err != nil {
				rb("submission attachments")
				util.HandleError(
					fmt.Errorf("error inserting multiple submission attachments in insert fetch_all data (chunk %d): %w", i, err),
				)
				return
			}
		}

		err = trx.Commit()
		if err != nil {
			rb("commit")
			util.HandleError(fmt.Errorf("error commiting insert fetch_all data transaction: %w", err))
			return
		}

		// success
		return
	}()
}
