package gradesapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/canvas_tokens"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/grades"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/notifications"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/sessions"
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
		s, err := session.NewSession(&aws.Config{Region: aws.String("us-east-2")})
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

	gradesIncludeSession        = gradesInclude("session")
	gradesIncludeUserProfile    = gradesInclude("user_profile")
	gradesIncludeObservees      = gradesInclude("observees")
	gradesIncludeCourses        = gradesInclude("courses")
	gradesIncludeOutcomeResults = gradesInclude("outcome_results")
	gradesIncludeSimpleGrades   = gradesInclude("simple_grades")
	gradesIncludeDetailedGrades = gradesInclude("detailed_grades")
	gradesIncludeGPA            = gradesInclude("gpa")
)

type gradesHandlerRequest struct {
	Session        bool
	UserProfile    bool
	Observees      bool
	Courses        bool
	OutcomeResults bool
	DetailedGrades bool
	GPA            bool
}

// UserGradesRequest represents a request for GradesForUser.
type UserGradesRequest struct {
	UserID         uint64
	CanvasUserID   uint64
	DetailedGrades bool
	ManualFetch    bool
}

// UserGradesResponse is all possible info from a GradesForUser call.
// It is JSON-serializable.
type UserGradesResponse struct {
	Session        *sessions.VerifiedSession `json:"session,omitempty"`
	UserProfile    *canvasUserProfile        `json:"user_profile,omitempty"`
	Observees      *[]canvasObservee         `json:"observees,omitempty"`
	Courses        *[]canvasCourse           `json:"courses,omitempty"`
	OutcomeResults processedOutcomeResults   `json:"outcome_results,omitempty"`
	SimpleGrades   simpleGrades              `json:"simple_grades,omitempty"`
	DetailedGrades detailedGrades            `json:"detailed_grades,omitempty"`
	GPA            gpa                       `json:"gpa"`
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
		s = append(s, oauth2.ScopeGrades)
	}

	if r.GPA {
		s = append(s, oauth2.ScopeGPA)
	}

	return s
}

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
			req.GPA = true
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

	g, gep := GradesForUser(&UserGradesRequest{
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
	toks, err := canvas_tokens.List(util.DB, &canvas_tokens.ListRequest{
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
		mutex = sync.Mutex{}
		wg    = sync.WaitGroup{}
		// error mapped by canvas user id
		errs = make(map[uint64]*GradesErrorResponse)
		// whether we had a success for the specified canvas user id
		// if 123 worked, statuses[123] = true.
		statuses = make(map[uint64]bool)
	)
	for _, tok := range *toks {
		wg.Add(1)
		go func(cuID uint64) {
			defer wg.Done()

			g, gep := GradesForUser(&UserGradesRequest{
				CanvasUserID: tok.CanvasUserID,
				// using DetailedGrades because it's computationally easier
				DetailedGrades: true,
				// not manual because this is a fetch for other users
				ManualFetch: false,
			})
			if gep != nil {
				mutex.Lock()

				errs[cuID] = gep
				statuses[cuID] = false

				mutex.Unlock()
				return
			}

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

			// list notification settings
			ns, err := notifications.ListSettings(db, &notifications.ListSettingsRequest{
				CanvasUserID: cuID,
				Type:         notifications.TypeGradeChange,
				Medium:       notifications.MediumEmail,
			})
			if err != nil {
				util.HandleError(fmt.Errorf("error listing notification settings for user %d: %w", cuID, err))
				set(false)
				return
			}

			if len(*ns) != 1 {
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

			before := time.Now().Add(-(2 * time.Minute))
			pgsP, err := grades.List(db, &grades.ListRequest{
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
		time.Sleep(150 * time.Millisecond)
	}

	wg.Wait()

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

func GradesForUser(req *UserGradesRequest) (*UserGradesResponse, *GradesErrorResponse) {
	var (
		rd  requestDetails
		err error
	)

	if req.UserID > 0 {
		rd, err = rdFromUserID(req.UserID)
	} else {
		rd, err = rdFromCanvasUserID(req.CanvasUserID)
	}

	if err != nil {
		return nil, &GradesErrorResponse{InternalError: fmt.Errorf("error getting rd from user id: %w", err)}
	}

	if rd.TokenID < 1 {
		return nil, &GradesErrorResponse{
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
					return nil, &GradesErrorResponse{
						Error:      gradesErrorRevokedToken,
						Action:     gradesErrorActionRedirectToOAuth,
						StatusCode: http.StatusForbidden,
					}
				}

				return nil, &GradesErrorResponse{
					InternalError: fmt.Errorf("error refreshing a token for a newProfile: %w", refreshErr),
				}
			}

			newProfile, newProfileErr := getCanvasProfile(rd, "self")
			if newProfileErr != nil {
				if errors.Is(newProfileErr, canvasErrorInvalidAccessTokenError) {
					return nil, &GradesErrorResponse{
						Error:      gradesErrorRefreshedTokenError,
						Action:     gradesErrorActionRedirectToOAuth,
						StatusCode: http.StatusForbidden,
					}
				} else if errors.Is(err, canvasErrorUnknownError) {
					return nil, &gradesErrorUnknownCanvasErrorResponse
				}

				return nil, &GradesErrorResponse{InternalError: fmt.Errorf("error getting a newProfile: %w", newProfileErr)}
			}

			profile = newProfile
		} else if errors.Is(err, canvasErrorInvalidAccessTokenError) {
			return nil, &GradesErrorResponse{
				Error:      gradesErrorRefreshedTokenError,
				Action:     gradesErrorActionRedirectToOAuth,
				StatusCode: http.StatusForbidden,
			}
		} else if errors.Is(err, canvasErrorUnknownError) {
			return nil, &gradesErrorUnknownCanvasErrorResponse
		} else {
			return nil, &GradesErrorResponse{InternalError: fmt.Errorf("error getting a canvas profile: %w", err)}
		}

		// reset err, this succeeded
		// in the future, err should always be nil
		err = nil
	}
	go saveProfileToDB((*canvasUserProfile)(profile))

	mutex := sync.Mutex{}
	wg := sync.WaitGroup{}

	var (
		allCourses *[]canvasCourse
		observees  *canvasUserObserveesResponse
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
			return nil, &GradesErrorResponse{
				Error:      gradesErrorRevokedToken,
				Action:     gradesErrorActionRetryOnce,
				StatusCode: http.StatusForbidden,
			}
		} else if errors.Is(err, canvasErrorUnknownError) {
			return nil, &gradesErrorUnknownCanvasErrorResponse
		}

		return nil, &GradesErrorResponse{InternalError: fmt.Errorf("error getting canvas courses: %w", err)}
	}

	go saveCoursesToDB(allCourses)
	go saveObserveesToDB((*[]canvasObservee)(observees), profile.ID)

	// we now have both allCourses and observees.
	gradedUsers, courses := getGradedUsersAndValidCourses(allCourses)

	// outcome_alignments / outcome_rollups / assignments [Grades/GradeBreakdown]

	// map[courseID]map[userID]map[outcomeID][]canvasOutcomeResult
	results := processedOutcomeResults{}

	for _, c := range *courses {
		// cID is a string of the course ID
		cID := strconv.Itoa(int(c.ID))

		// uIDs is a string slice of all graded users in the course
		var uIDs []string
		for _, uID := range gradedUsers[c.ID] {
			uIDs = append(uIDs, strconv.Itoa(int(uID)))
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

	// wait for data
	wg.Wait()

	if err != nil {
		if errors.Is(err, canvasErrorInvalidAccessTokenError) {
			return nil, &GradesErrorResponse{
				Error:         gradesErrorRevokedToken,
				Action:        gradesErrorActionRetryOnce,
				StatusCode:    http.StatusForbidden,
				InternalError: nil,
			}
		} else if errors.Is(err, canvasErrorUnknownError) {
			return nil, &gradesErrorUnknownCanvasErrorResponse
		}

		return nil, &GradesErrorResponse{InternalError: fmt.Errorf("error getting alignments/results/assignments: %w", err)}
	}

	go saveOutcomeResultsToDB(results)

	// now, we will calculate grades
	// map[userID<uint64>]map[courseID<uint64>]grade<computedGrade>
	grades := detailedGrades{}
	sGrades := simpleGrades{}

	for cID, uIDs := range gradedUsers {
		for _, uID := range uIDs {
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
					var c canvasCourse
					for _, co := range *courses {
						if co.ID == courseID {
							c = co
							break
						}
					}
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

	go saveGradesToDB(grades, req.ManualFetch)

	cGPA := calculateGPAFromDetailedGrades(grades)

	go saveGPAToDB(cGPA, req.ManualFetch)

	return &UserGradesResponse{
		Session:        nil,
		UserProfile:    (*canvasUserProfile)(profile),
		Observees:      (*[]canvasObservee)(observees),
		Courses:        courses,
		OutcomeResults: results,
		SimpleGrades:   sGrades,
		DetailedGrades: grades,
		GPA:            cGPA,
	}, nil
}
