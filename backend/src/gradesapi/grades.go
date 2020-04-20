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
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/enrollments"
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
	UserID       uint64
	CanvasUserID uint64
	// If specified, don't fetch grades for specified courses. Not respected in AllGradesForTeacher.
	ExcludeCourseIDs map[uint64]struct{}
	// Not respected in AllGradesForTeacher.
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
	Enrollments            []enrollments.UpsertRequest
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

	// 1. Figure out all courses to pull (from enrollments.List)

	// first, just teachers
	teacherEnrollments, err := enrollments.List(db, &enrollments.ListRequest{Type: enrollments.TypeTeacher})
	if err != nil {
		e := fmt.Errorf("error listing teacher enrollments in fetch_all: %w", err)
		util.HandleError(e)
		uploadToS3(true, e)
		return
	}

	// teachersByTeacher (map[teacherID<uint64>][]courseID<uint64>) tells us teachers and the classes they teach.
	teachersByTeacher := map[uint64][]uint64{}
	// teachersByCourse (map[courseID<uint64>]teacherID<uint64> quickly lets us figure out if a course has a teacher (that's signed up!) or not.
	teachersByCourse := map[uint64]uint64{}
	/* students (map[userID<uint64>][]courseID<uint64>) tells us which courses students are in that a teacher is
	available for.
	*/
	studentsExcludeCourses := map[uint64][]uint64{}
	// studentsAll is the opposite of studentsExcludeCourses.
	studentsAll := map[uint64][]uint64{}

	for _, e := range *teacherEnrollments {
		teachersByTeacher[e.UserCanvasID] = append(teachersByTeacher[e.UserCanvasID], e.CourseID)
		teachersByCourse[e.CourseID] = e.UserCanvasID
	}

	// now, list student enrollments.
	studentEnrollments, err := enrollments.List(db, &enrollments.ListRequest{Type: enrollments.TypeStudent})
	if err != nil {
		e := fmt.Errorf("error listing student enrollments in fetch_all: %w", err)
		util.HandleError(e)
		uploadToS3(true, e)
		return
	}

	for _, e := range *studentEnrollments {
		studentsAll[e.UserCanvasID] = append(studentsAll[e.UserCanvasID], e.CourseID)

		// if a teacher exists
		if _, ok := teachersByCourse[e.CourseID]; ok {
			studentsExcludeCourses[e.UserCanvasID] = append(studentsExcludeCourses[e.UserCanvasID], e.CourseID)
		}
	}

	// 2. Get all grade change email recipients and their previous grades.

	// students that want notifications (map[studentCanvasUserID<uint64>]struct{}{})
	studentsEnabledNotifications := make(map[uint64]struct{})
	var studentsEnabledNotificationsSlice []uint64

	notificationReqs, err := notifications.ListSettings(db, &notifications.ListSettingsRequest{
		Type:   notifications.TypeGradeChange,
		Medium: notifications.MediumEmail,
	})
	if err != nil {
		e := fmt.Errorf("error listing notification requests in fetch_all: %w", err)
		util.HandleError(e)
		uploadToS3(true, e)
		return
	}

	for _, r := range *notificationReqs {
		studentsEnabledNotifications[r.CanvasUserID] = struct{}{}
		studentsEnabledNotificationsSlice = append(studentsEnabledNotificationsSlice, r.CanvasUserID)
	}

	// for those, we'll get their previous grades
	prevGrades, err := gradessvc.List(db, &gradessvc.ListRequest{
		UserCanvasIDs: &studentsEnabledNotificationsSlice,
	})
	if err != nil {
		e := fmt.Errorf("error listing previous grades in fetch_all: %w", err)
		util.HandleError(e)
		uploadToS3(true, e)
		return
	}

	// holds all previous grades for students who have notifications enabled
	// map[studentID<uint64>]map[courseID<uint64>]gradessvc.Grade
	studentPrevGrades := make(map[uint64]map[uint64]gradessvc.Grade, len(studentsEnabledNotifications))
	// holds all new grades for students who have notifications enabled
	studentNewGrades := make(map[uint64]map[uint64]computedGrade, len(studentsEnabledNotifications))
	for _, pg := range *prevGrades {
		if studentPrevGrades[pg.UserCanvasID] == nil {
			studentPrevGrades[pg.UserCanvasID] = map[uint64]gradessvc.Grade{pg.CourseID: pg}
		} else {
			studentPrevGrades[pg.UserCanvasID][pg.CourseID] = pg
		}
	}

	// 3. Grab tokens for everyone.

	tokens, err := canvas_tokens.List(util.DB, &canvas_tokens.ListRequest{
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

	var (
		teacherTokens []canvas_tokens.CanvasToken
		// bothTokens covers a VERY RARE circumstance that a user is both a student and a teacher in different courses.
		bothTokens     []canvas_tokens.CanvasToken
		studentTokens  []canvas_tokens.CanvasToken
		observerTokens []canvas_tokens.CanvasToken
	)

	for _, t := range *tokens {
		// if the user is a teacher
		_, userIsTeacher := teachersByTeacher[t.CanvasUserID]
		_, userIsStudent := studentsAll[t.CanvasUserID]

		if userIsTeacher && userIsStudent {
			bothTokens = append(bothTokens, t)
			continue
		} else if userIsTeacher && !userIsStudent {
			teacherTokens = append(teacherTokens, t)
			continue
		} else if !userIsTeacher && userIsStudent {
			studentTokens = append(studentTokens, t)
		} else if !userIsTeacher && !userIsStudent {
			// user is an observer because they are not a student and not a teacher
			observerTokens = append(observerTokens, t)
		}
	}

	// 4. If so, fetch teacher first, then students (without teacher courses). If not, fetch everything.

	// error mapped by canvas user id
	errs := make(map[uint64]*GradesErrorResponse)
	// whether we had a success for the specified canvas user id
	// if 123 worked, statuses[123] = true.
	teacherStatuses := make(map[uint64]bool)
	// excludeCourses contains courses that a teacher has already fetched the grade for.
	excludeCourses := map[uint64]struct{}{}
	// map[courseID]courseName
	excludedCourseNames := map[uint64]string{}

	var dbReqs []UserGradesDBRequests

	// concatenate teachers and both
	for _, tt := range append(teacherTokens, bothTokens...) {
		rd := rdFromToken(tt)

		resp, dbReq, err := AllGradesForTeacher(&UserGradesRequest{
			CanvasUserID:     tt.CanvasUserID,
			ManualFetch:      false,
			ReturnDBRequests: true,
			Rd:               &rd,
		})
		if err != nil {
			if err.InternalError != nil {
				util.HandleError(fmt.Errorf("error in fetch_all when getching grades for teacher %d: %w", tt.CanvasUserID, err.InternalError))
			}

			errs[tt.CanvasUserID] = err
		} else {
			dbReqs = append(dbReqs, *dbReq)
		}

		teacherStatuses[tt.CanvasUserID] = err == nil

		if err != nil {
			continue
		}

		for _, c := range *resp.Courses {
			excludeCourses[c.ID] = struct{}{}
			excludedCourseNames[c.ID] = c.Name
		}

		// store grades for students with notifications enabled
		for uID, dg := range resp.DetailedGrades {
			if _, ok := studentsEnabledNotifications[uID]; ok {
				if studentNewGrades[uID] == nil {
					studentNewGrades[uID] = map[uint64]computedGrade{}
				}

				for cID, g := range dg {
					studentNewGrades[uID][cID] = g
				}
			}
		}
	}

	// to DB we go
	if len(dbReqs) > 0 {
		handleBatchGradesDBRequests(dbReqs)
	}

	// students, both, observers

	dbReqs = []UserGradesDBRequests{}
	restStatuses := make(map[uint64]bool)

	handleRestRequests := func(t canvas_tokens.CanvasToken) {
		rd := rdFromToken(t)

		resp, dbReq, err := GradesForUser(&UserGradesRequest{
			CanvasUserID:     t.CanvasUserID,
			ExcludeCourseIDs: excludeCourses,
			DetailedGrades:   true,
			ManualFetch:      false,
			ReturnDBRequests: true,
			Rd:               &rd,
		})
		if err != nil {
			if err.InternalError != nil {
				util.HandleError(fmt.Errorf("error in fetch_all when getching grades for user %d: %w", t.CanvasUserID, err.InternalError))
			}

			errs[t.CanvasUserID] = err
		} else {
			dbReqs = append(dbReqs, *dbReq)
		}

		restStatuses[t.CanvasUserID] = err == nil

		if err != nil {
			return
		}

		userIsObserver := resp.Observees != nil && len(*resp.Observees) > 0

		if _, ok := studentsEnabledNotifications[t.CanvasUserID]; ok {
			// the user would like a notification

			courseNames := make(map[uint64]string)
			for _, c := range *resp.Courses {
				courseNames[c.ID] = c.Name
			}

			pr := *resp.UserProfile

			// once for grades a teacher fetched
			for cID, g := range studentNewGrades[pr.ID] {
				// if a previous grade exists for user
				if uPrev, ok := studentPrevGrades[pr.ID]; ok {
					// prev grade for course?
					if prev, ok := uPrev[cID]; ok {
						// are they different?
						if g.Grade.Grade != prev.Grade {
							courseName := courseNames[cID]
							if len(courseName) < 1 {
								courseName = excludedCourseNames[cID]
							}

							if userIsObserver {
								var studentName string
								for _, o := range *resp.Observees {
									if o.ID == pr.ID {
										studentName = o.Name
										break
									}
								}

								go email.SendParentGradeChangeEmail(&email.ParentGradeChangeEmailData{
									To:            pr.PrimaryEmail,
									Name:          pr.Name,
									StudentName:   studentName,
									ClassName:     courseName,
									PreviousGrade: prev.Grade,
									CurrentGrade:  g.Grade.Grade,
								})
							} else {
								go email.SendGradeChangeEmail(&email.GradeChangeEmailData{
									To:            pr.PrimaryEmail,
									Name:          pr.Name,
									ClassName:     courseName,
									PreviousGrade: prev.Grade,
									CurrentGrade:  g.Grade.Grade,
								})
							}
						}
					}
				}
			}

			// and once more for grades the user fetched
			for uID, cs := range resp.DetailedGrades {
				// range thru courses
				for cID, c := range cs {
					// if a previous grade exists for user
					if uPrev, ok := studentPrevGrades[uID]; ok {
						// prev grade for course?
						if prev, ok := uPrev[cID]; ok {
							// are they different?
							if c.Grade.Grade != prev.Grade {
								if userIsObserver {
									var studentName string
									for _, o := range *resp.Observees {
										if o.ID == uID {
											studentName = o.Name
											break
										}
									}

									go email.SendParentGradeChangeEmail(&email.ParentGradeChangeEmailData{
										To:            pr.PrimaryEmail,
										Name:          pr.Name,
										StudentName:   studentName,
										ClassName:     courseNames[cID],
										PreviousGrade: prev.Grade,
										CurrentGrade:  c.Grade.Grade,
									})
								} else {
									go email.SendGradeChangeEmail(&email.GradeChangeEmailData{
										To:            pr.PrimaryEmail,
										Name:          pr.Name,
										ClassName:     courseNames[cID],
										PreviousGrade: prev.Grade,
										CurrentGrade:  c.Grade.Grade,
									})
								}
							}
						}
					}
				}
			}
		}
	}

	// splitting these due to "ON CONFLICT can't update same row twice" in enrollments
	for _, st := range append(studentTokens, bothTokens...) {
		handleRestRequests(st)
	}

	// to DB we go
	if len(dbReqs) > 0 {
		handleBatchGradesDBRequests(dbReqs)
	}

	dbReqs = []UserGradesDBRequests{}

	for _, ot := range observerTokens {
		handleRestRequests(ot)
	}

	// to DB we go
	if len(dbReqs) > 0 {
		handleBatchGradesDBRequests(dbReqs)
	}

	resp := struct {
		Errors          map[uint64]*GradesErrorResponse `json:"errors"`
		TeacherStatuses map[uint64]bool                 `json:"teacher_statuses"`
		RestStatuses    map[uint64]bool                 `json:"rest_statuses"`
		NumTeachers     int                             `json:"num_teachers"`
		NumBoth         int                             `json:"num_both"`
		NumStudents     int                             `json:"num_students"`
		NumObservers    int                             `json:"num_observers"`
		NumErrors       int                             `json:"num_errors"`
	}{
		Errors:          errs,
		TeacherStatuses: teacherStatuses,
		RestStatuses:    restStatuses,
		NumTeachers:     len(teacherTokens),
		NumBoth:         len(bothTokens),
		NumStudents:     len(studentTokens),
		NumObservers:    len(observerTokens),
		NumErrors:       len(errs),
	}

	if returnData {
		jRet, err := json.Marshal(&resp)
		if err != nil {
			handleISE(w, fmt.Errorf("error marshaling errors and statuses from fetch all grades: %w", err))
			return
		}

		util.SendJSONResponse(w, jRet)
		return
	} else {
		uploadToS3(false, resp)
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

			cReq, eReq := prepareCoursesForDB(allCourses)
			dbReqs.Courses = cReq
			dbReqs.Enrollments = *eReq
		}()
	} else {
		go saveCoursesToDB(allCourses)
	}

	// observees are special and don't get added to dbReqs
	go saveObserveesToDB((*[]canvasObservee)(observees), profile.ID)

	// we now have both allCourses and observees.
	gradedUsers, validCourses := getGradedUsersAndValidCourses(allCourses)
	var courses []canvasCourse

	if req.ExcludeCourseIDs != nil {
		for _, c := range *validCourses {
			if _, ok := req.ExcludeCourseIDs[c.ID]; !ok {
				courses = append(courses, c)
			}
		}
	} else {
		courses = *validCourses
	}

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

// AllGradesForTeacher fetches grades for all of a teacher's classes and all of said classes' students.
func AllGradesForTeacher(req *UserGradesRequest) (*UserGradesResponse, *UserGradesDBRequests, *GradesErrorResponse) {
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
		allCourses *[]canvasCourse
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
				for _, e := range c.Enrollments {
					// make sure that the user is a teacher in this course
					if e.Type == enrollments.TypeTeacher {
						cs = append(cs, c)
						break
					}
				}
			}
		}

		allCourses = &cs
		mutex.Unlock()
		return
	}()

	// wait
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

			cReq, eReq := prepareCoursesForDB(allCourses)
			dbReqs.Courses = cReq
			dbReqs.Enrollments = *eReq
		}()
	} else {
		go saveCoursesToDB(allCourses)
	}

	// map[courseID]map[userID]map[outcomeID][]canvasOutcomeResult
	results := processedOutcomeResults{}
	enrolls := make(map[uint64][]canvasFullEnrollment, len(*allCourses))
	var (
		allEnrolls  []canvasFullEnrollment
		dlCourseIDs []uint64
	)

	// get outcome results
	for _, c := range *allCourses {
		if c.EnrollmentTermID != spring20DLEnrollmentTermID {
			// results
			wg.Add(1)
			go func(courseID uint64) {
				defer wg.Done()

				rs, rErr := getCanvasOutcomeResults(
					rd,
					fmt.Sprintf("%d", courseID),
					// we want for all!
					[]string{},
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
			}(c.ID)
		} else {
			dlCourseIDs = append(dlCourseIDs, c.ID)
		}

		// fetch enrollments
		wg.Add(1)
		go func(courseID uint64) {
			defer wg.Done()

			es, eErr := getCanvasCourseEnrollments(
				rd,
				fmt.Sprintf("%d", courseID),
			)
			if eErr != nil {
				mutex.Lock()
				err = eErr
				mutex.Unlock()
				return
			}

			mutex.Lock()
			enrolls[courseID] = *es
			allEnrolls = append(allEnrolls, *es...)
			mutex.Unlock()
		}(c.ID)
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

	// TODO: Send teachers emails when calls fail
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

		dbReqsWg.Add(1)
		go func() {
			defer dbReqsWg.Done()

			dbReqs.Enrollments = append(dbReqs.Enrollments, *prepareEnrollmentsForDB(allEnrolls)...)
		}()
	} else {
		go saveOutcomeResultsToDB(results)
		go saveEnrollmentsToDB(allEnrolls)
	}

	// calculate grades

	grades := detailedGrades{}

	// first, distance learning
	for _, cID := range dlCourseIDs {
		for _, e := range enrolls[cID] {
			if len(e.Grades.CurrentGrade) < 1 {
				continue
			}

			if grades[e.UserID] == nil {
				grades[e.UserID] = map[uint64]computedGrade{
					e.CourseID: {Grade: grade{Grade: e.Grades.CurrentGrade}}}
			} else {
				grades[e.UserID][e.CourseID] = computedGrade{
					Grade: grade{Grade: e.Grades.CurrentGrade},
				}
			}
		}
	}

	// second, CBL

	for cID, us := range results {
		for uID, rs := range us {
			wg.Add(1)
			go func(courseID uint64, userID uint64, scores map[uint64][]canvasOutcomeResult) {
				defer wg.Done()

				// we're saying it's not after the cutoff for now.
				grd := *calculateGradeFromOutcomeResults(scores, false)

				// we'll now save the grade
				mutex.Lock()

				if grades[userID] == nil {
					grades[userID] = make(map[uint64]computedGrade)
				}

				grades[userID][courseID] = grd
				mutex.Unlock()
				return
			}(cID, uID, rs)
		}
	}

	// wait for CBL
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
		Courses:        allCourses,
		OutcomeResults: results,
		DetailedGrades: grades,
		//GPA:            cGPA,
		DistanceLearning: dlGrades,
	}, &dbReqs, nil
}
