package gradesapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/sessions"
	"github.com/iamtheyammer/canvascbl/backend/src/middlewares"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"sync"
)

var (
	db                                    = util.DB
	gradesErrorUnknownCanvasErrorResponse = gradesErrorResponse{
		Error: gradesErrorUnknownCanvasError,
	}
)

type gradesErrorAction string
type gradesInclude string

// map[courseTitle<string>]map[userID<uint64>]grade<string>
type simpleGrades map[string]map[uint64]string

// map[userID<uint64>]map[courseID<uint64>]grade<computedGrade>
type detailedGrades map[uint64]map[uint64]computedGrade

// map[courseID]map[userID]map[outcomeID][]canvasOutcomeResult
type processedOutcomeResults map[uint64]map[uint64]map[uint64][]canvasOutcomeResult

const (
	gradesErrorNoTokens              = "no stored tokens for this user"
	gradesErrorRevokedToken          = "the token/refresh token has been revoked or no longer works"
	gradesErrorRefreshedTokenError   = "after refreshing the token, it is invalid"
	gradesErrorUnknownCanvasError    = "there was an unknown error from canvas"
	gradesErrorInvalidInclude        = "invalid include"
	gradesErrorActionRedirectToOAuth = gradesErrorAction("redirect_to_oauth")
	gradesErrorActionRetryOnce       = gradesErrorAction("retry_once")

	gradesIncludeSession        = gradesInclude("session")
	gradesIncludeUserProfile    = gradesInclude("user_profile")
	gradesIncludeObservees      = gradesInclude("observees")
	gradesIncludeCourses        = gradesInclude("courses")
	gradesIncludeOutcomeResults = gradesInclude("outcome_results")
	gradesIncludeSimpleGrades   = gradesInclude("simple_grades")
	gradesIncludeDetailedGrades = gradesInclude("detailed_grades")
)

type gradesHandlerRequest struct {
	Session        bool
	UserProfile    bool
	Observees      bool
	Courses        bool
	OutcomeResults bool
	DetailedGrades bool
}

type gradesHandlerResponse struct {
	Session        *sessions.VerifiedSession `json:"session,omitempty"`
	UserProfile    *canvasUserProfile        `json:"user_profile,omitempty"`
	Observees      *[]canvasObservee         `json:"observees,omitempty"`
	Courses        *[]canvasCourse           `json:"courses,omitempty"`
	OutcomeResults processedOutcomeResults   `json:"outcome_results,omitempty"`
	SimpleGrades   simpleGrades              `json:"simple_grades,omitempty"`
	DetailedGrades detailedGrades            `json:"detailed_grades,omitempty"`
}

type gradesErrorResponse struct {
	Error  string            `json:"error"`
	Action gradesErrorAction `json:"action,omitempty"`
}

func GradesHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// user (this is where we'll handle the refresh token)
	// allCourses / observees [Dashboard]
	// outcome_alignments / outcome_rollups / outcome_results / assignments [Grades/GradeBreakdown]
	// individual outcomes [GradeBreakdown]
	session := middlewares.Session(w, r)
	if session == nil {
		return
	}

	req := gradesHandlerRequest{}
	if session.Type == sessions.VerifiedSessionTypeSessionString {
		inc := r.URL.Query()["include[]"]

		// everything is permitted, what do they want?
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
			default:
				handleError(w, gradesErrorResponse{
					Error: gradesErrorInvalidInclude,
				}, http.StatusBadRequest)
				return
			}
		}
	} else if session.Type == sessions.VerifiedSessionTypeAPIKey {
		util.SendUnauthorized(w, "api keys aren't implemented yet")
		return
	} else {
		util.SendUnauthorized(w, "unsupported authentication method")
		return
	}

	rd, err := rdFromCanvasUserID(session.CanvasUserID)
	if err != nil {
		handleISE(w, fmt.Errorf("error getting rd from canvas user id: %w", err))
		return
	}

	if rd.TokenID < 1 {
		handleError(w, gradesErrorResponse{
			Error:  gradesErrorNoTokens,
			Action: gradesErrorActionRedirectToOAuth,
		}, http.StatusForbidden)
		return
	}

	profile, err := getCanvasProfile(rd, "self")
	if err != nil {
		if errors.Is(err, canvasErrorInvalidAccessTokenError) {
			// we need to use the refresh token
			refreshErr := rd.refreshAccessToken()
			if refreshErr != nil {
				if errors.Is(refreshErr, canvasErrorInvalidAccessTokenError) {
					handleError(w, gradesErrorResponse{
						Error:  gradesErrorRevokedToken,
						Action: gradesErrorActionRedirectToOAuth,
					}, http.StatusForbidden)
					return
				}

				handleISE(w, fmt.Errorf("error getting a newProfile: %w", refreshErr))
				return
			}

			newProfile, newProfileErr := getCanvasProfile(rd, "self")
			if newProfileErr != nil {
				if errors.Is(newProfileErr, canvasErrorInvalidAccessTokenError) {
					handleError(w, gradesErrorResponse{
						Error:  gradesErrorRefreshedTokenError,
						Action: gradesErrorActionRedirectToOAuth,
					}, http.StatusForbidden)
					return
				} else if errors.Is(err, canvasErrorUnknownError) {
					handleError(w, gradesErrorUnknownCanvasErrorResponse, util.CanvasProxyErrorCode)
					return
				}

				handleISE(w, fmt.Errorf("error getting a newProfile: %w", newProfileErr))
				return
			}

			profile = newProfile
		} else if errors.Is(err, canvasErrorInvalidAccessTokenError) {
			handleError(w, gradesErrorResponse{
				Error:  gradesErrorRefreshedTokenError,
				Action: gradesErrorActionRedirectToOAuth,
			}, http.StatusForbidden)
			return
		} else if errors.Is(err, canvasErrorUnknownError) {
			handleError(w, gradesErrorUnknownCanvasErrorResponse, util.CanvasProxyErrorCode)
			return
		} else {
			handleISE(w, fmt.Errorf("error getting a canvas profile: %w", err))
			return
		}

		// reset err, this succeeded
		// in the future, err should always be nil
		err = nil
	}
	go saveProfileToDB((*canvasUserProfile)(profile))

	mutex := sync.Mutex{}
	wg := sync.WaitGroup{}

	var (
		allCourses *canvasCoursesResponse
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

		allCourses = coursesResp
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
			handleError(w, gradesErrorResponse{
				Error:  gradesErrorRevokedToken,
				Action: gradesErrorActionRetryOnce,
			}, http.StatusForbidden)
			return
		} else if errors.Is(err, canvasErrorUnknownError) {
			handleError(w, gradesErrorUnknownCanvasErrorResponse, util.CanvasProxyErrorCode)
			return
		}

		handleISE(w, fmt.Errorf("error getting canvas courses: %w", err))
		return
	}

	go saveCoursesToDB((*[]canvasCourse)(allCourses))
	go saveObserveesToDB((*[]canvasObservee)(observees), profile.ID)

	// we now have both allCourses and observees.
	gradedUsers, courses := getGradedUsersAndValidCourses((*[]canvasCourse)(allCourses))

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
			handleError(w, gradesErrorResponse{
				Error:  gradesErrorRevokedToken,
				Action: gradesErrorActionRetryOnce,
			}, http.StatusForbidden)
			return
		} else if errors.Is(err, canvasErrorUnknownError) {
			handleError(w, gradesErrorUnknownCanvasErrorResponse, util.CanvasProxyErrorCode)
			return
		}

		handleISE(w, fmt.Errorf("error getting alignments/results/assignments: %w", err))
		return
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

	go saveGradesToDB(grades)

	resp := gradesHandlerResponse{}

	if req.Session {
		resp.Session = session
	}

	if req.UserProfile {
		resp.UserProfile = (*canvasUserProfile)(profile)
	}

	if req.Observees {
		resp.Observees = (*[]canvasObservee)(observees)
	}

	if req.Courses {
		resp.Courses = courses
	}

	if req.OutcomeResults {
		resp.OutcomeResults = results
	}

	if req.DetailedGrades {
		resp.DetailedGrades = grades
	} else {
		resp.SimpleGrades = sGrades
	}

	jResp, err := json.Marshal(&resp)
	if err != nil {
		handleISE(w, fmt.Errorf("error marshaling grades handler response into JSON: %w", err))
		return
	}
	util.SendJSONResponse(w, jResp)

	return
}
