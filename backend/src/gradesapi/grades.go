package gradesapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	awssession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/canvas_tokens"
	coursessvc "github.com/iamtheyammer/canvascbl/backend/src/db/services/courses"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/gpas"
	gradessvc "github.com/iamtheyammer/canvascbl/backend/src/db/services/grades"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/notifications"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/sessions"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/users"
	"github.com/iamtheyammer/canvascbl/backend/src/email"
	"github.com/iamtheyammer/canvascbl/backend/src/env"
	"github.com/iamtheyammer/canvascbl/backend/src/middlewares"
	"github.com/iamtheyammer/canvascbl/backend/src/oauth2"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var (
	db                                    = util.DB
	gradesErrorUnknownCanvasErrorResponse = GradesErrorResponse{
		Error:      gradesErrorUnknownCanvasError,
		StatusCode: util.CanvasProxyErrorCode,
	}
	s3Uploader = func() *s3manager.Uploader {
		s, err := awssession.NewSession(&aws.Config{Region: aws.String("us-east-2")})
		if err != nil {
			panic(fmt.Errorf("error creating aws session: %w", err))
		}

		ul := s3manager.NewUploader(s)

		return ul
	}()
)

type gradesErrorAction string
type gradesInclude string

// map[courseTitle<string>]map[userID<uint64>]grade<string>
type simpleGrades map[string]map[uint64]string

// map[userID<uint64>]map[courseID<uint64>]grade<computedGrade>
type detailedGrades map[uint64]map[uint64]computedGrade

// map[courseID]map[userID]map[outcomeID][]canvasOutcomeResult
type processedOutcomeResults map[uint64]map[uint64]map[uint64][]canvasOutcomeResult

// calculatedGPA represents a single user's gpa
type calculatedGPA struct {
	Unweighted struct {
		Subgrades float64 `json:"subgrades"`
		Default   float64 `json:"default"`
	} `json:"unweighted"`
}

// gpa represents more than one user's GPA
type gpa map[uint64]calculatedGPA

type distanceLearningGrade struct {
	CourseName string `json:"course_name"`
	Grade      struct {
		Grade string `json:"grade"`
		Rank  int    `json:"rank"`
	} `json:"grade"`
	OriginalCourseID         uint64 `json:"original_course_id"`
	DistanceLearningCourseID uint64 `json:"distance_learning_course_id"`
}

// map[userID<uint64>][]distanceLearningGrade
type distanceLearningGrades map[uint64][]distanceLearningGrade

const (
	gradesErrorNoTokens              = "no stored tokens for this user"
	gradesErrorRevokedToken          = "the token/refresh token has been revoked or no longer works"
	gradesErrorRefreshedTokenError   = "after refreshing the token, it is invalid"
	gradesErrorUnknownCanvasError    = "there was an unknown error from canvas"
	gradesErrorInvalidInclude        = "invalid include"
	gradesErrorUnauthorizedScope     = "your oauth2 grant doesn't have one or more requested scopes"
	gradesErrorInvalidAccessToken    = "invalid access token"
	gradesErrorActionRedirectToOAuth = gradesErrorAction("redirect_to_oauth")
	gradesErrorActionRetryOnce       = gradesErrorAction("retry_once")

	gradesIncludeSession          = gradesInclude("session")
	gradesIncludeUserProfile      = gradesInclude("user_profile")
	gradesIncludeObservees        = gradesInclude("observees")
	gradesIncludeCourses          = gradesInclude("courses")
	gradesIncludeOutcomeResults   = gradesInclude("outcome_results")
	gradesIncludeSimpleGrades     = gradesInclude("simple_grades")
	gradesIncludeDetailedGrades   = gradesInclude("detailed_grades")
	gradesIncludeGPA              = gradesInclude("gpa")
	gradesIncludeDistanceLearning = gradesInclude("distance_learning")
)

type gradesHandlerRequest struct {
	Session          bool
	UserProfile      bool
	Observees        bool
	Courses          bool
	OutcomeResults   bool
	DetailedGrades   bool
	GPA              bool
	DistanceLearning bool
}

// UserGradesRequest represents a request for GradesForUser.
type UserGradesRequest struct {
	UserID           uint64
	CanvasUserID     uint64
	DetailedGrades   bool
	ManualFetch      bool
	ReturnDBRequests bool
	Rd               *requestDetails
}

// UserGradesResponse is all possible info from a GradesForUser call.
// It is JSON-serializable.
type UserGradesResponse struct {
	Session          *sessions.VerifiedSession `json:"session,omitempty"`
	UserProfile      *canvasUserProfile        `json:"user_profile,omitempty"`
	Observees        *[]canvasObservee         `json:"observees,omitempty"`
	Courses          *[]canvasCourse           `json:"courses,omitempty"`
	OutcomeResults   processedOutcomeResults   `json:"outcome_results,omitempty"`
	SimpleGrades     simpleGrades              `json:"simple_grades,omitempty"`
	DetailedGrades   detailedGrades            `json:"detailed_grades,omitempty"`
	GPA              gpa                       `json:"gpa,omitempty"`
	DistanceLearning distanceLearningGrades    `json:"distance_learning,omitempty"`
}

/*
UserGradesDBRequests is all of the insert/upsert requests that would be
performed during the execution of this function.

It is returned from GradesForUser if req.ReturnDBRequests is true.

Note that observees are excluded due to their special upsert nature.
*/
type UserGradesDBRequests struct {
	Profile                *users.UpsertRequest
	Courses                *[]coursessvc.UpsertRequest
	OutcomeResults         *[]coursessvc.OutcomeResultInsertRequest
	Grades                 *[]gradessvc.InsertRequest
	RollupScores           *[]coursessvc.OutcomeRollupInsertRequest
	GPA                    *[]gpas.InsertRequest
	DistanceLearningGrades *[]gradessvc.InsertDistanceLearningRequest
}

// GradesErrorResponse represents an error from GradesForUser.
// InternalError will be populated when there is a server error.
// It is JSON-serializable.
type GradesErrorResponse struct {
	Error         string            `json:"error"`
	Action        gradesErrorAction `json:"action,omitempty"`
	StatusCode    int               `json:"status_code,omitempty"`
	InternalError error             `json:"-"`
}

func (r gradesHandlerRequest) toScopes() []oauth2.Scope {
	var s []oauth2.Scope

	gradesScope := false

	// session not supported

	if r.UserProfile {
		s = append(s, oauth2.ScopeProfile)
	}

	if r.Observees {
		s = append(s, oauth2.ScopeObservees)
	}

	if r.Courses {
		s = append(s, oauth2.ScopeCourses)
	}

	if r.OutcomeResults {
		s = append(s, oauth2.ScopeOutcomeResults)
	}

	if r.DetailedGrades {
		s = append(s, oauth2.ScopeDetailedGrades)
	} else {
		gradesScope = true
		s = append(s, oauth2.ScopeGrades)
	}

	if r.DistanceLearning && !gradesScope {
		gradesScope = true
		s = append(s, oauth2.ScopeGrades)
	}

	if r.GPA {
		s = append(s, oauth2.ScopeGPA)
	}

	return s
}

// GradesHandler handles /api/v1/grades
func GradesHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	inc := r.URL.Query()["include[]"]
	req := gradesHandlerRequest{}

	for _, i := range inc {
		switch gradesInclude(i) {
		case gradesIncludeSession:
			req.Session = true
		case gradesIncludeUserProfile:
			req.UserProfile = true
		case gradesIncludeObservees:
			req.Observees = true
		case gradesIncludeCourses:
			req.Courses = true
		case gradesIncludeOutcomeResults:
			req.OutcomeResults = true
		case gradesIncludeSimpleGrades:
		case gradesIncludeDetailedGrades:
			req.DetailedGrades = true
		case gradesIncludeGPA:
		// Distance Learning
		//req.GPA = true
		case gradesIncludeDistanceLearning:
			req.DistanceLearning = true
		default:
			handleError(w, GradesErrorResponse{
				Error: gradesErrorInvalidInclude,
			}, http.StatusBadRequest)
			return
		}
	}

	var (
		at, tokenIsOK = middlewares.Bearer(w, r, false)
		session       *sessions.VerifiedSession
		userID        uint64
	)

	if !tokenIsOK {
		handleError(w, GradesErrorResponse{
			Error: gradesErrorInvalidAccessToken,
		}, http.StatusUnauthorized)
		return
	}

	if len(at) < 1 {
		// session time
		session = middlewares.Session(w, r, true)
		if session == nil {
			return
		}

		userID = session.UserID
	} else {
		// oauth2
		if req.Session {
			// invalid
			handleError(w, GradesErrorResponse{
				Error: gradesErrorInvalidInclude,
			}, http.StatusBadRequest)
			return
		}
		grant, err := oauth2.Authorizer(at, req.toScopes(), &oauth2.AuthorizerAPICall{
			RoutePath: "grades",
			Method:    "GET",
			Query:     &r.URL.RawQuery,
		})
		if err != nil {
			if errors.Is(err, oauth2.GrantMissingScopeError) {
				handleError(w, GradesErrorResponse{
					Error: gradesErrorUnauthorizedScope,
				}, http.StatusUnauthorized)
				return
			}

			if errors.Is(err, oauth2.InvalidAccessTokenError) {
				handleError(w, GradesErrorResponse{
					Error: oauth2.InvalidAccessTokenError.Error(),
				}, http.StatusForbidden)
				return
			}

			handleISE(w, fmt.Errorf("error using oauth2.Authorizer in GradesHandler: %w", err))
			return
		}

		userID = grant.UserID
	}

	// a fetch is considered manual if it's initiated with a session.
	// it also has to be via this endpoint, so that's already covered.
	manualFetch := false
	if session != nil {
		manualFetch = true
	}

	g, _, gep := GradesForUser(&UserGradesRequest{
		UserID:         userID,
		DetailedGrades: req.DetailedGrades,
		ManualFetch:    manualFetch,
	})
	if gep != nil {
		if gep.InternalError != nil {
			handleISE(w, gep.InternalError)
			return
		}
		handleError(w, *gep, gep.StatusCode)
		return
	}

	resp := UserGradesResponse{}

	if req.Session {
		resp.Session = session
	}

	if req.UserProfile {
		resp.UserProfile = g.UserProfile
	}

	if req.Observees {
		resp.Observees = g.Observees
	}

	if req.Courses {
		resp.Courses = g.Courses
	}

	if req.OutcomeResults {
		resp.OutcomeResults = g.OutcomeResults
	}

	if req.DetailedGrades {
		resp.DetailedGrades = g.DetailedGrades
	} else {
		resp.SimpleGrades = g.SimpleGrades
	}

	if req.GPA {
		resp.GPA = g.GPA
	}

	if req.DistanceLearning {
		resp.DistanceLearning = g.DistanceLearning
	}

	jResp, err := json.Marshal(&resp)
	if err != nil {
		handleISE(w, fmt.Errorf("error marshaling grades handler response into JSON: %w", err))
		return
	}
	util.SendJSONResponse(w, jResp)

	return
}

// GradesForAllHandler gets grades for all users with a specified key.
func GradesForAllHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	key := r.Header.Get("X-CanvasCBL-Script-Key")
	if len(key) < 1 {
		util.SendBadRequest(w, "missing X-CanvasCBL-Script-Key as header")
		return
	} else if key != env.ScriptKey {
		util.SendUnauthorized(w, "invalid X-CanvasCBL-Script-Key as header")
		return
	}

	returnData := r.URL.Query().Get("return_data") == "true"
	if !returnData {
		util.SendNoContent(w)
		// NOT returning-- we want to finish our work
	}

	uploadToS3 := func(isError bool, input interface{}) {
		key := time.Now().Format(time.RFC3339) + "-" + string(env.Env) + ".json"
		if isError {
			key = "error-" + key
		}

		i := &s3manager.UploadInput{
			Bucket:      aws.String("canvascbl-fetch-all-grades-logs"),
			ContentType: aws.String("application/json"),
			Key:         aws.String(key),
		}

		jRet, err := json.Marshal(input)
		if err != nil {
			util.HandleError(fmt.Errorf("error marshaling upload to s3 input to json: %w", err))
		}

		i.Body = bytes.NewReader(jRet)

		_, err = s3Uploader.Upload(i)
		if err != nil {
			util.HandleError(fmt.Errorf("error uploading to s3: %w", err))
		}
	}

	// get all users with tokens
	toksP, err := canvas_tokens.List(util.DB, &canvas_tokens.ListRequest{
		OrderBys:   []string{"canvas_tokens.canvas_user_id", "canvas_tokens.inserted_at DESC"},
		DistinctOn: "canvas_tokens.canvas_user_id",
	})
	if err != nil {
		wErr := fmt.Errorf("error listing all unique canvas tokens for grades for all: %w", err)
		if returnData {
			handleISE(w, wErr)
		} else {
			uploadToS3(true, &wErr)
		}

		return
	}

	toks := *toksP

	var (
		mutex = sync.Mutex{}
		wg    = sync.WaitGroup{}
		// error mapped by canvas user id
		errs = make(map[uint64]*GradesErrorResponse)
		// whether we had a success for the specified canvas user id
		// if 123 worked, statuses[123] = true.
		statuses = make(map[uint64]bool)

		dbReqsMutex = sync.Mutex{}
		dbReqsWg    = sync.WaitGroup{}
		dbReqs      []*UserGradesDBRequests
	)

	notificationReqs, err := notifications.ListSettings(db, &notifications.ListSettingsRequest{
		Type:   notifications.TypeGradeChange,
		Medium: notifications.MediumEmail,
	})
	if err != nil {
		util.HandleError(fmt.Errorf("error listing notification requests in fetch_all handler: %w", err))
		return
	}

	// notification requests, by canvas user id
	nrs := make(map[uint64]struct{})
	for _, ns := range *notificationReqs {
		nrs[ns.CanvasUserID] = struct{}{}
	}

	delay := time.Duration(len(toks)) * time.Millisecond

	for _, tok := range toks {
		wg.Add(1)
		go func(cuID uint64) {
			defer wg.Done()

			rd := rdFromToken(tok)

			g, dbReq, gep := GradesForUser(&UserGradesRequest{
				CanvasUserID: tok.CanvasUserID,
				// using DetailedGrades because it's computationally easier
				DetailedGrades: true,
				// not manual because this is a fetch for other users
				ManualFetch: false,
				// we literally already fetched the token
				Rd: &rd,
				// batch em!
				ReturnDBRequests: true,
			})
			if gep != nil {
				mutex.Lock()

				errs[cuID] = gep
				statuses[cuID] = false

				mutex.Unlock()
				return
			}

			dbReqsWg.Add(1)
			go func() {
				defer dbReqsWg.Done()

				dbReqsMutex.Lock()
				dbReqs = append(dbReqs, dbReq)
				dbReqsMutex.Unlock()
			}()

			/*
				now, we'll handle a grade change notification
			*/

			// sets the status of the request, to be used before returning
			set := func(ok bool) {
				mutex.Lock()
				statuses[cuID] = ok
				mutex.Unlock()
			}

			if len(g.DetailedGrades) < 1 {
				set(true)
				return
			}

			if _, ok := nrs[tok.CanvasUserID]; !ok {
				set(true)
				return
			}

			// now, we know that they do want an email.

			var (
				userIDs   []uint64
				courseIDs = make(map[uint64]struct{})
			)
			for uID, cs := range g.DetailedGrades {
				// user has one or more classes with a grade?
				validUser := false
				for cID, grd := range cs {
					if grd.Grade != naGrade {
						validUser = true

						courseIDs[cID] = struct{}{}
					}
				}

				if validUser {
					userIDs = append(userIDs, uID)
				}
			}
			if len(userIDs) < 1 {
				return
			}

			var cIDs []uint64
			for cID := range courseIDs {
				cIDs = append(cIDs, cID)
			}

			// TODO: remove database call
			before := time.Now().Add(-(2 * time.Minute))
			pgsP, err := gradessvc.List(db, &gradessvc.ListRequest{
				UserCanvasIDs: &userIDs,
				Before:        &before,
				CourseIDs:     &cIDs,
			})
			if err != nil {
				util.HandleError(fmt.Errorf("error listing previous grades for user %d: %w", cuID, err))
				set(false)
				return
			}
			pgs := *pgsP

			type change struct {
				CanvasUserID  uint64
				CourseID      uint64
				PreviousGrade string
				CurrentGrade  string
			}

			var (
				fetchUserIDs   = make(map[uint64]struct{})
				fetchCourseIDs = make(map[uint64]struct{})
				changes        []change
			)

			// check if grades have changed.

			// check all users
			for uID, cs := range g.DetailedGrades {
				// and all courses
				for cID, grd := range cs {
					// if the new grade is N/A, skip course
					if grd.Grade == naGrade {
						continue
					}

					// then loop thru previous grades for a match
					for _, pg := range pgs {
						if pg.UserCanvasID == uID && pg.CourseID == cID && pg.Grade != grd.Grade.Grade {
							fetchUserIDs[uID] = struct{}{}
							fetchCourseIDs[cID] = struct{}{}
							changes = append(changes, change{CanvasUserID: uID, CourseID: cID, PreviousGrade: pg.Grade, CurrentGrade: grd.Grade.Grade})
							break
						}
					}
				}
			}

			courseNames := make(map[uint64]string)
			for _, c := range *g.Courses {
				courseNames[c.ID] = c.Name
			}

			userIsParent := len(*g.Observees) > 0
			for _, c := range changes {
				p := *g.UserProfile
				if userIsParent {
					go func(pr canvasUserProfile, ch change, obs []canvasObservee, courseName string) {
						// we need to search for the observee's name
						var studentName string
						for _, o := range obs {
							if o.ID == ch.CanvasUserID {
								studentName = o.Name
								break
							}
						}

						email.SendParentGradeChangeEmail(&email.ParentGradeChangeEmailData{
							To:            pr.PrimaryEmail,
							Name:          pr.Name,
							StudentName:   studentName,
							ClassName:     courseName,
							PreviousGrade: ch.PreviousGrade,
							CurrentGrade:  ch.CurrentGrade,
						})
					}(p, c, *g.Observees, courseNames[c.CourseID])
				} else {
					go email.SendGradeChangeEmail(&email.GradeChangeEmailData{
						To:            p.PrimaryEmail,
						Name:          p.Name,
						ClassName:     courseNames[c.CourseID],
						PreviousGrade: c.PreviousGrade,
						CurrentGrade:  c.CurrentGrade,
					})
				}
			}

			set(true)
			return
		}(tok.CanvasUserID)
		time.Sleep(delay)
	}

	wg.Wait()
	dbReqsWg.Wait()

	go func() {
		var (
			profiles []users.UpsertRequest
			// keeps out duplicates
			coursesMap = make(map[int64]struct{})
			courses    []coursessvc.UpsertRequest
			// keeps out duplicates
			outcomeResultsMap = make(map[uint64]struct{})
			// chunked in 7281 due to postgres's 65535 parameter limit
			// 9 params each
			chunkedOutcomeResults           = [][]coursessvc.OutcomeResultInsertRequest{{}}
			currentOutcomeResultChunk       = 0
			currentOutcomeResultChunkLength = 0
			// chunk in 16383 due to postgres's 65535 parameter limit
			// 4 params each
			chunkedGrades            = [][]gradessvc.InsertRequest{{}}
			currentGradesChunk       = 0
			currentGradesChunkLength = 0
			// chunk in 13107 due to postgres's 65535 parameter limit
			// 5 params each
			chunkedRollupScores           = [][]coursessvc.OutcomeRollupInsertRequest{{}}
			currentRollupScoreChunk       = 0
			currentRollupScoreChunkLength = 0
			gpaReqs                       []gpas.InsertRequest
			distanceLearningGradesReqs    []gradessvc.InsertDistanceLearningRequest
		)

		for _, r := range dbReqs {
			if r.Profile != nil {
				profiles = append(profiles, *r.Profile)
			}

			if r.Courses != nil {
				// keeps out duplicate courses
				for _, c := range *r.Courses {
					if _, ok := coursesMap[c.CourseID]; !ok {
						courses = append(courses, c)
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
							chunkedOutcomeResults = append(chunkedOutcomeResults, []coursessvc.OutcomeResultInsertRequest{})
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
						chunkedGrades = append(chunkedGrades, []gradessvc.InsertRequest{})
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
						chunkedRollupScores = append(chunkedRollupScores, []coursessvc.OutcomeRollupInsertRequest{})
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

		err = coursessvc.UpsertMultiple(trx, &courses)
		if err != nil {
			rb("courses")
			util.HandleError(fmt.Errorf("error upserting multiple courses in insert fetch_all data: %w", err))
			return
		}

		for i, req := range chunkedOutcomeResults {
			err = coursessvc.InsertMultipleOutcomeResults(trx, &req)
			if err != nil {
				rb("outcome results")
				util.HandleError(
					fmt.Errorf("error inserting multiple outcome results in insert fetch_all data (chunk %d): %w", i, err),
				)
				return
			}
		}

		for i, req := range chunkedGrades {
			err = gradessvc.Insert(trx, &req)
			if err != nil {
				rb("grades")
				util.HandleError(
					fmt.Errorf("error inserting multiple grades in insert fetch_all data (chunk %d): %w", i, err),
				)
				return
			}
		}

		for i, req := range chunkedRollupScores {
			err = coursessvc.InsertMultipleOutcomeRollups(trx, &req)
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

		err = gradessvc.InsertDistanceLearning(trx, &distanceLearningGradesReqs)
		if err != nil {
			rb("distance learning grades")
			util.HandleError(
				fmt.Errorf("error inserting distance learning grades in insert fetch_all data: %w", err),
			)
			return
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

	if returnData {
		jRet, err := json.Marshal(&struct {
			Errors   map[uint64]*GradesErrorResponse `json:"errors"`
			Statuses map[uint64]bool                 `json:"statuses"`
		}{
			Errors:   errs,
			Statuses: statuses,
		})
		if err != nil {
			handleISE(w, fmt.Errorf("error marshaling errors and statuses from fetch all grades: %w", err))
			return
		}

		util.SendJSONResponse(w, jRet)
		return
	} else {
		uploadToS3(false, &struct {
			Errors   map[uint64]*GradesErrorResponse `json:"errors"`
			Statuses map[uint64]bool                 `json:"statuses"`
		}{
			Errors:   errs,
			Statuses: statuses,
		})
	}

	return
}

func GradesForUser(req *UserGradesRequest) (*UserGradesResponse, *UserGradesDBRequests, *GradesErrorResponse) {
	var (
		rd       requestDetails
		dbReqs   UserGradesDBRequests
		dbReqsWg = sync.WaitGroup{}
		err      error
	)

	if req.Rd == nil {
		if req.UserID > 0 {
			rd, err = rdFromUserID(req.UserID)
		} else {
			rd, err = rdFromCanvasUserID(req.CanvasUserID)
		}
	} else {
		rd = *req.Rd
	}

	if err != nil {
		return nil, nil, &GradesErrorResponse{InternalError: fmt.Errorf("error getting rd from user id: %w", err)}
	}

	if rd.TokenID < 1 {
		return nil, nil, &GradesErrorResponse{
			Error:      gradesErrorNoTokens,
			Action:     gradesErrorActionRedirectToOAuth,
			StatusCode: http.StatusForbidden,
		}
	}

	profile, err := getCanvasProfile(rd, "self")
	if err != nil {
		if errors.Is(err, canvasErrorInvalidAccessTokenError) {
			// we need to use the refresh token
			refreshErr := rd.refreshAccessToken()
			if refreshErr != nil {
				if errors.Is(refreshErr, canvasErrorInvalidAccessTokenError) ||
					errors.Is(refreshErr, canvasOAuth2ErrorRefreshTokenNotFound) {
					return nil, nil, &GradesErrorResponse{
						Error:      gradesErrorRevokedToken,
						Action:     gradesErrorActionRedirectToOAuth,
						StatusCode: http.StatusForbidden,
					}
				}

				return nil, nil, &GradesErrorResponse{
					InternalError: fmt.Errorf("error refreshing a token for a newProfile: %w", refreshErr),
				}
			}

			newProfile, newProfileErr := getCanvasProfile(rd, "self")
			if newProfileErr != nil {
				if errors.Is(newProfileErr, canvasErrorInvalidAccessTokenError) {
					return nil, nil, &GradesErrorResponse{
						Error:      gradesErrorRefreshedTokenError,
						Action:     gradesErrorActionRedirectToOAuth,
						StatusCode: http.StatusForbidden,
					}
				} else if errors.Is(err, canvasErrorUnknownError) {
					return nil, nil, &gradesErrorUnknownCanvasErrorResponse
				}

				return nil, nil, &GradesErrorResponse{InternalError: fmt.Errorf("error getting a newProfile: %w", newProfileErr)}
			}

			profile = newProfile
		} else if errors.Is(err, canvasErrorInvalidAccessTokenError) {
			return nil, nil, &GradesErrorResponse{
				Error:      gradesErrorRefreshedTokenError,
				Action:     gradesErrorActionRedirectToOAuth,
				StatusCode: http.StatusForbidden,
			}
		} else if errors.Is(err, canvasErrorUnknownError) {
			return nil, nil, &gradesErrorUnknownCanvasErrorResponse
		} else {
			return nil, nil, &GradesErrorResponse{InternalError: fmt.Errorf("error getting a canvas profile: %w", err)}
		}

		// reset err, this succeeded
		// in the future, err should always be nil
		err = nil
	}

	if req.ReturnDBRequests {
		dbReqsWg.Add(1)
		go func() {
			defer dbReqsWg.Done()

			dbReqs.Profile = prepareProfileForDB((*canvasUserProfile)(profile))
		}()
	} else {
		go saveProfileToDB((*canvasUserProfile)(profile))
	}

	mutex := sync.Mutex{}
	wg := sync.WaitGroup{}

	var (
		allCourses    *[]canvasCourse
		hiddenCourses map[uint64]struct{}
		observees     *canvasUserObserveesResponse
	)

	// get allCourses
	wg.Add(1)
	go func() {
		defer wg.Done()

		coursesResp, coursesErr := getCanvasCourses(rd)
		mutex.Lock()
		if coursesErr != nil {
			err = coursesErr
			mutex.Unlock()
			return
		}

		var cs []canvasCourse
		for _, c := range *coursesResp {
			if int(c.EnrollmentTermID) >= env.CanvasCurrentEnrollmentTermID {
				cs = append(cs, c)
			}
		}

		allCourses = &cs
		mutex.Unlock()
		return
	}()

	// we don't want to bother with this if we're bunching them all
	if !req.ReturnDBRequests {
		// get the user's hidden courses
		wg.Add(1)
		go func() {
			defer wg.Done()

			hiddenIDs, hiddenErr := coursessvc.GetUserHiddenCourses(db, req.UserID)
			if hiddenErr != nil {
				mutex.Lock()
				err = hiddenErr
				mutex.Unlock()
				return
			}

			hiddenCourses = *hiddenIDs
			return
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		observeesResp, observeesErr := getCanvasUserObservees(rd, "self")
		mutex.Lock()
		if observeesErr != nil {
			err = observeesErr
			mutex.Unlock()
			return
		}

		observees = observeesResp
		mutex.Unlock()
		return
	}()

	// wait for both to finish
	wg.Wait()

	if err != nil {
		if errors.Is(err, canvasErrorInvalidAccessTokenError) {
			return nil, nil, &GradesErrorResponse{
				Error:      gradesErrorRevokedToken,
				Action:     gradesErrorActionRetryOnce,
				StatusCode: http.StatusForbidden,
			}
		} else if errors.Is(err, canvasErrorUnknownError) {
			return nil, nil, &gradesErrorUnknownCanvasErrorResponse
		}

		return nil, nil, &GradesErrorResponse{InternalError: fmt.Errorf("error getting canvas courses: %w", err)}
	}

	if req.ReturnDBRequests {
		dbReqsWg.Add(1)
		go func() {
			defer dbReqsWg.Done()

			dbReqs.Courses = prepareCoursesForDB(allCourses)
		}()
	} else {
		go saveCoursesToDB(allCourses)
	}

	// observees are special and don't get added to dbReqs
	go saveObserveesToDB((*[]canvasObservee)(observees), profile.ID)

	// we now have both allCourses and observees.
	gradedUsers, coursesP := getGradedUsersAndValidCourses(allCourses)
	courses := *coursesP

	// outcome_alignments / outcome_rollups / assignments [Grades/GradeBreakdown]

	// map[courseID]map[userID]map[outcomeID][]canvasOutcomeResult
	results := processedOutcomeResults{}

	for i, c := range courses {
		if c.EnrollmentTermID != spring20DLEnrollmentTermID {
			// cID is a string of the course ID
			cID := strconv.Itoa(int(c.ID))

			// uIDs is a string slice of all graded users in the course
			var uIDs []string
			for _, uID := range gradedUsers[c.ID] {
				uIDs = append(uIDs, strconv.Itoa(int(uID)))
			}

			if _, ok := hiddenCourses[c.ID]; ok {
				courses[i].CanvasCBLHidden = true
			}

			// results
			wg.Add(1)
			go func(courseIDS string, courseID uint64) {
				defer wg.Done()

				rs, rErr := getCanvasOutcomeResults(
					rd,
					courseIDS,
					uIDs,
				)
				if rErr != nil {
					mutex.Lock()
					err = rErr
					mutex.Unlock()
					return
				}

				processedResults, processErr := processOutcomeResults(&rs.OutcomeResults)
				if processErr != nil {
					mutex.Lock()
					err = processErr
					mutex.Unlock()
					return
				}

				mutex.Lock()
				results[courseID] = *processedResults
				mutex.Unlock()
				return
			}(cID, c.ID)
		}
	}

	// wait for data
	wg.Wait()

	if err != nil {
		if errors.Is(err, canvasErrorInvalidAccessTokenError) {
			return nil, nil, &GradesErrorResponse{
				Error:         gradesErrorRevokedToken,
				Action:        gradesErrorActionRetryOnce,
				StatusCode:    http.StatusForbidden,
				InternalError: nil,
			}
		} else if errors.Is(err, canvasErrorUnknownError) {
			return nil, nil, &gradesErrorUnknownCanvasErrorResponse
		}

		return nil, nil, &GradesErrorResponse{InternalError: fmt.Errorf("error getting alignments/results/assignments: %w", err)}
	}

	if req.ReturnDBRequests {
		dbReqsWg.Add(1)
		go func() {
			defer dbReqsWg.Done()

			resReq, prepErr := prepareOutcomeResultsForDB(results)
			if prepErr != nil {
				util.HandleError(fmt.Errorf("error preparing outcome results for db: %w", err))
				return
			}

			dbReqs.OutcomeResults = resReq
		}()
	} else {
		go saveOutcomeResultsToDB(results)
	}

	// now, we will calculate grades
	// map[userID<uint64>]map[courseID<uint64>]grade<computedGrade>
	grades := detailedGrades{}
	sGrades := simpleGrades{}

	for cID, uIDs := range gradedUsers {
		// course object
		var c canvasCourse
		for _, cc := range *allCourses {
			if cc.ID == cID {
				c = cc
				break
			}
		}

		for _, uID := range uIDs {
			if c.EnrollmentTermID == spring20DLEnrollmentTermID {
				// we will be using Canvas's grade, found in the enrollment object.
				var e canvasEnrollment
				for _, en := range c.Enrollments {
					if en.UserID == uID {
						e = en
						break
					}
				}

				// if a computed current grade exists
				if len(e.ComputedCurrentGrade) > 0 {
					if !req.DetailedGrades {
						if sGrades[c.Name] == nil {
							sGrades[c.Name] = make(map[uint64]string)
						}

						// drop it into simple grades (easy!)
						sGrades[c.Name][uID] = e.ComputedCurrentGrade
					} else {
						if grades[uID] == nil {
							grades[uID] = make(map[uint64]computedGrade)
						}

						grades[uID][cID] = computedGrade{
							// just make a new grade object.
							Grade: grade{Grade: e.ComputedCurrentGrade},
							// so we get [] instead of null
							//Averages: make(map[uint64]computedAverage),
						}
					}
				} else {
					if !req.DetailedGrades {
						if sGrades[c.Name] == nil {
							sGrades[c.Name] = make(map[uint64]string)
						}
						sGrades[c.Name][uID] = naGrade.Grade
					}

					if grades[uID] == nil {
						grades[uID] = make(map[uint64]computedGrade)
					}

					grades[uID][cID] = computedGrade{
						Grade: naGrade,
						// so we get [] instead of null
						Averages: make(map[uint64]computedAverage),
					}
				}

				continue
			}
			wg.Add(1)
			go func(courseID uint64, userID uint64) {
				defer wg.Done()

				mutex.Lock()
				rs := results[courseID][userID]
				mutex.Unlock()

				// we're saying it's not after the cutoff for now.
				grd := *calculateGradeFromOutcomeResults(rs, false)

				// we'll now save the grade
				mutex.Lock()
				if !req.DetailedGrades {
					if sGrades[c.Name] == nil {
						sGrades[c.Name] = make(map[uint64]string)
					}
					sGrades[c.Name][userID] = grd.Grade.Grade
				}

				if grades[userID] == nil {
					grades[userID] = make(map[uint64]computedGrade)
				}

				grades[userID][courseID] = grd
				mutex.Unlock()
				return
			}(cID, uID)
		}
	}

	wg.Wait()

	if req.ReturnDBRequests {
		dbReqsWg.Add(1)
		go func() {
			defer dbReqsWg.Done()

			gReqs, rollupReqs := prepareGradesForDB(grades, req.ManualFetch)

			dbReqs.Grades = gReqs
			dbReqs.RollupScores = rollupReqs
		}()
	} else {
		go saveGradesToDB(grades, req.ManualFetch)
	}

	dlGrades := distanceLearningGrades{}

	for userID, dg := range grades {
		wg.Add(1)
		go func(uID uint64, ac []canvasCourse, detGra map[uint64]computedGrade) {
			defer wg.Done()

			//fmt.Printf("detailed grades for user %d: %+v\n\n", uID, detGra)
			dlg := calculateDistanceLearningGrades(ac, detGra)

			mutex.Lock()
			dlGrades[uID] = dlg
			mutex.Unlock()
		}(userID, *allCourses, dg)
	}

	wg.Wait()

	if req.ReturnDBRequests {
		dbReqsWg.Add(1)
		go func() {
			defer dbReqsWg.Done()

			dbReqs.DistanceLearningGrades = prepareDistanceLearningGradesForDB(dlGrades, req.ManualFetch)
		}()
	} else {
		go saveDistanceLearningGradesToDB(dlGrades, req.ManualFetch)
	}

	//cGPA := calculateGPAFromDetailedGrades(grades)
	//
	//if req.ReturnDBRequests {
	//	dbReqsWg.Add(1)
	//	go func() {
	//		defer dbReqsWg.Done()
	//
	//		dbReqs.GPA = prepareGPAForDB(cGPA, req.ManualFetch)
	//	}()
	//
	//	// wait for them all to finish
	//	dbReqsWg.Wait()
	//} else {
	//	go saveGPAToDB(cGPA, req.ManualFetch)
	//}

	if req.ReturnDBRequests {
		dbReqsWg.Wait()
	}

	return &UserGradesResponse{
		Session:        nil,
		UserProfile:    (*canvasUserProfile)(profile),
		Observees:      (*[]canvasObservee)(observees),
		Courses:        &courses,
		OutcomeResults: results,
		SimpleGrades:   sGrades,
		DetailedGrades: grades,
		//GPA:            cGPA,
		DistanceLearning: dlGrades,
	}, &dbReqs, nil
}
