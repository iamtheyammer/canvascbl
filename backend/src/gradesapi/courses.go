package gradesapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/enrollments"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/submissions"
	userssvc "github.com/iamtheyammer/canvascbl/backend/src/db/services/users"
	"github.com/iamtheyammer/canvascbl/backend/src/oauth2"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

type listCoursesResponse struct {
	Courses               *canvasCoursesResponse  `json:"courses"`
	DistanceLearningPairs *[]distanceLearningPair `json:"distance_learning_pairs,omitempty"`
}

type listEnrollmentsResponse struct {
	Enrollments []canvasFullEnrollment `json:"enrollments"`
}

type lateSubmissionSummary struct {
	Total   uint64  `json:"total"`
	Percent float64 `json:"percent"`

	Graded        uint64 `json:"graded"`
	PendingReview uint64 `json:"pending_review"`
	Submitted     uint64 `json:"submitted"`
	Unsubmitted   uint64 `json:"unsubmitted"`
}

type submissionSummary struct {
	Total         uint64                 `json:"total"`
	Graded        uint64                 `json:"graded"`
	PendingReview uint64                 `json:"pending_review"`
	Submitted     uint64                 `json:"submitted"`
	Unsubmitted   uint64                 `json:"unsubmitted"`
	Late          *lateSubmissionSummary `json:"late,omitempty"`
}

type courseUserSubmissionSummaryResponse struct {
	SubmissionSummary map[uint64]submissionSummary `json:"submission_summary"`
}

// ListCoursesHandler lists courses for a user.
func ListCoursesHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userID, rdP, sess := authorizer(w, r, []oauth2.Scope{oauth2.ScopeCourses}, &oauth2.AuthorizerAPICall{
		Method:    "GET",
		RoutePath: "courses",
		Query:     &r.URL.RawQuery,
	})
	if (userID == nil || rdP == nil) && sess == nil {
		return
	}

	returnDLPairs := false
	if r.URL.Query().Get("include[]") == "distance_learning_pairs" {
		returnDLPairs = true
	}

	var cs *canvasCoursesResponse
	_, err := handleRequestWithTokenRefresh(func(reqD *requestDetails) error {
		tCourses, tErr := getCanvasCourses(*reqD)
		if tErr != nil {
			return fmt.Errorf("error getting canvas courses: %w", tErr)
		}

		cs = tCourses
		return nil
	}, rdP, *userID)
	if errors.Is(err, canvasErrorInvalidAccessTokenError) || errors.Is(err, canvasErrorInsufficientScopesOnAccessTokenError) {
		handleError(w, GradesErrorResponse{
			Error:  gradesErrorRevokedToken,
			Action: gradesErrorActionRedirectToOAuth,
		}, http.StatusForbidden)
		return
	} else if errors.Is(err, canvasErrorUnknownError) {
		handleError(w, gradesErrorUnknownCanvasErrorResponse, util.CanvasProxyErrorCode)
		return
	} else if err != nil {
		handleISE(w, fmt.Errorf("error getting courses: %w", err))
		return
	}

	go saveCoursesToDB((*[]canvasCourse)(cs))

	var dlPairs *[]distanceLearningPair
	if returnDLPairs {
		dlPairs = findDistanceLearningCoursePairs(*cs)
	}

	j, err := json.Marshal(&listCoursesResponse{
		Courses:               cs,
		DistanceLearningPairs: dlPairs,
	})
	if err != nil {
		handleISE(w, fmt.Errorf("error marshaling list courses response to json: %w", err))
		return
	}

	util.SendJSONResponse(w, j)
}

// CourseEnrollmentsHandler lists enrollments for a course. /api/v1/courses/:courseID/enrollments
func CourseEnrollmentsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	courseID := ps.ByName("courseID")
	if len(courseID) < 1 || !util.ValidateIntegerString(courseID) {
		util.SendBadRequest(w, "missing or invalid courseID as url param")
		return
	}

	q := r.URL.Query()
	types := q["type[]"]
	for _, t := range types {
		switch t {
		case "StudentEnrollment":
		case "TeacherEnrollment":
		case "TaEnrollment":
		case "DesignerEnrollment":
		case "ObserverEnrollment":
		default:
			util.SendBadRequest(w, "invalid type as query param")
			return
		}
	}

	states := q["state[]"]
	for _, s := range states {
		switch s {
		case "active":
		case "invited":
		case "creation_pending":
		case "deleted":
		case "rejected":
		case "completed":
		case "inactive":
		case "current_and_invited":
		case "current_and_future":
		case "current_and_concluded":
		default:
			util.SendBadRequest(w, "invalid state as query param")
			return
		}
	}

	userID, rdP, sess := authorizer(w, r, []oauth2.Scope{oauth2.ScopeEnrollments, oauth2.ScopeGrades}, &oauth2.AuthorizerAPICall{
		Method:    "GET",
		RoutePath: "courses/:courseID/enrollments",
		Query:     &r.URL.RawQuery,
	})
	if (userID == nil || rdP == nil) && sess == nil {
		return
	}

	var es canvasEnrollmentsResponse
	_, err := handleRequestWithTokenRefresh(func(reqD *requestDetails) error {
		tEnrollments, tErr := getCanvasCourseEnrollments(*reqD, &getCanvasCourseEnrollmentsRequest{
			courseID: courseID,
			types:    types,
			states:   states,
			includes: []string{"avatar_url", "observed_users"},
		})
		if tErr != nil {
			return fmt.Errorf("error getting canvas course enrollments for course %s: %w", courseID, tErr)
		}

		es = *tEnrollments
		return nil
	}, rdP, *userID)
	if errors.Is(err, canvasErrorInvalidAccessTokenError) || errors.Is(err, canvasErrorInsufficientScopesOnAccessTokenError) {
		handleError(w, GradesErrorResponse{
			Error:  gradesErrorRevokedToken,
			Action: gradesErrorActionRedirectToOAuth,
		}, http.StatusForbidden)
		return
	} else if errors.Is(err, canvasErrorUnknownError) {
		handleError(w, gradesErrorUnknownCanvasErrorResponse, util.CanvasProxyErrorCode)
		return
	} else if err != nil {
		handleISE(w, fmt.Errorf("error getting courses: %w", err))
		return
	}

	go saveEnrollmentsToDB(es)

	sendJSON(w, &listEnrollmentsResponse{Enrollments: es})
	return
}

//CourseSubmissionSummaryHandler handles /api/v1/courses/:courseID/submission_summary/users
func CourseSubmissionSummaryHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// params
	courseID := ps.ByName("courseID")
	if len(courseID) < 1 || !util.ValidateIntegerString(courseID) {
		util.SendBadRequest(w, "missing or invalid courseID as url param")
		return
	}

	cID, err := strconv.Atoi(courseID)
	if err != nil {
		util.SendBadRequest(w, "invalid courseID as url param")
		return
	}

	q := r.URL.Query()

	userIDs := q["user_ids[]"]
	var uIDs []uint64
	uIDsMap := make(map[uint64]struct{}, len(userIDs))
	if len(userIDs) > 0 {
		for _, strID := range userIDs {
			id, err := strconv.Atoi(strID)
			if err != nil || id < 1 {
				util.SendBadRequest(w, "invalid user_id as query param")
				return
			}
			if _, ok := uIDsMap[uint64(id)]; !ok {
				uIDs = append(uIDs, uint64(id))
				uIDsMap[uint64(id)] = struct{}{}
			}
		}
	}

	useCache := q.Get("use_cache") == "true"

	var lateCount bool

	for _, incl := range q["include[]"] {
		switch incl {
		case "late_count":
			lateCount = true
		default:
			util.SendBadRequest(w, "invalid include as query param")
			return
		}
	}

	// authorizer

	userID, rdP, sess := authorizer(w, r, []oauth2.Scope{oauth2.ScopeSubmissions}, &oauth2.AuthorizerAPICall{
		Method:    "GET",
		RoutePath: "courses/:courseID/submission_summary/users",
		Query:     &r.URL.RawQuery,
	})
	if (userID == nil || rdP == nil) && sess == nil {
		return
	}

	rd := *rdP

	if !useCache && !rd.hasScopeVersion(2) {
		handleError(w, GradesErrorResponse{
			Error:      gradesErrorMissingCanvasScope,
			Action:     gradesErrorActionRedirectToOAuth,
			StatusCode: http.StatusUnauthorized,
		}, http.StatusUnauthorized)
		return
	}

	users, err := userssvc.List(db, &userssvc.ListRequest{ID: *userID})
	if err != nil {
		handleISE(w, fmt.Errorf("error getting calling user in submission summary: %w", err))
		return
	}

	if len(*users) < 1 {
		handleISE(w, fmt.Errorf("couldn't find calling user in submission summary with user ID %d", *userID))
		return
	}

	user := (*users)[0]

	// permissions

	// first, we need to know the calling user's enrollment.
	var callingUserEnrollmentType enrollments.Type

	// get from Canvas
	var ces canvasEnrollmentsResponse
	rd, err = handleRequestWithTokenRefresh(func(reqD *requestDetails) error {
		tResp, outErr := getCanvasCourseEnrollments(*reqD, &getCanvasCourseEnrollmentsRequest{
			courseID: courseID,
			userID:   "self",
			states:   []string{"active"},
		})
		if outErr != nil {
			return fmt.Errorf("error getting assignments for course %s: %w", cID, outErr)
		}

		ces = *tResp
		return nil
	}, &rd, *userID)
	if errors.Is(err, canvasErrorInvalidAccessTokenError) || errors.Is(err, canvasErrorInsufficientScopesOnAccessTokenError) {
		handleError(w, GradesErrorResponse{
			Error:  gradesErrorRevokedToken,
			Action: gradesErrorActionRedirectToOAuth,
		}, http.StatusForbidden)
		return
	} else if errors.Is(err, canvasErrorUnknownError) {
		handleError(w, gradesErrorUnknownCanvasErrorResponse, util.CanvasProxyErrorCode)
		return
	} else if err != nil {
		handleISE(w, fmt.Errorf("error getting canvas course enrollments in submission summary: %w", err))
		return
	}

	// convert to enrollments.Enrollment
	var types []enrollments.Type
	for _, ce := range ces {
		types = append(types, enrollments.Role(ce.Role).ToType())
	}

	callingUserEnrollmentType = enrollments.MostPermissiveType(types...)

	if !callingUserEnrollmentType.Valid() {
		util.SendUnauthorized(w, "unable to find a valid enrollment for you in this course")
		return
	}

	if !callingUserEnrollmentType.OneOf(enrollments.TypeTeacher, enrollments.TypeStudent, enrollments.TypeObserver) {
		util.SendUnauthorized(w, "your enrollment type in this course does not support submission summaries at this time")
		return
	}

	// get allowed users
	var (
		permittedUserIDsMap = make(map[uint64]struct{})
		userSpecifiedIDs    = len(uIDs) > 0
		requestedUserIDs    = uIDs
	)

	if callingUserEnrollmentType == enrollments.TypeTeacher {
		// nothing
	} else if callingUserEnrollmentType == enrollments.TypeStudent {
		permittedUserIDsMap[user.CanvasUserID] = struct{}{}
		if !userSpecifiedIDs {
			requestedUserIDs = append(requestedUserIDs, user.CanvasUserID)
		}
	} else if callingUserEnrollmentType == enrollments.TypeObserver {
		// must go fetch observees
		obs, err := userssvc.ListObservees(db, &userssvc.ListObserveesRequest{
			ObserverCanvasUserID: user.CanvasUserID,
			ActiveOnly:           true,
		})
		if err != nil {
			handleISE(w, fmt.Errorf("error fetching observees in submission summary: %w", err))
			return
		}

		for _, o := range *obs {
			permittedUserIDsMap[o.CanvasUserID] = struct{}{}

			if !userSpecifiedIDs {
				requestedUserIDs = append(requestedUserIDs, o.CanvasUserID)
			}
		}
	}

	// teachers may get all data
	if callingUserEnrollmentType != enrollments.TypeTeacher {
		for id := range uIDsMap {
			if _, ok := permittedUserIDsMap[id]; !ok {
				util.SendUnauthorized(w, "you do not have access to one or more of the requested user's submissions")
				return
			}
		}
	}

	// force cache for entire class
	if callingUserEnrollmentType == enrollments.TypeTeacher && len(uIDsMap) < 1 {
		useCache = true
	}

	// initialize response
	resp := courseUserSubmissionSummaryResponse{
		SubmissionSummary: map[uint64]submissionSummary{},
	}

	if useCache {
		sums, err := submissions.GetCourseUserSummary(db, &submissions.CourseUserSummaryRequest{
			CourseID:     uint64(cID),
			UserIDs:      requestedUserIDs,
			SeparateLate: lateCount,
		})
		if err != nil {
			handleISE(w, fmt.Errorf("error getting course user submission summary from db for"+
				" course %d and user(s) %v: %w", cID, uIDs, err))
			return
		}

		for uID, dbSumm := range *sums {
			summ := submissionSummary{
				Graded:        dbSumm.Graded,
				PendingReview: dbSumm.PendingReview,
				Submitted:     dbSumm.Submitted,
				Unsubmitted:   dbSumm.Unsubmitted,
			}

			var l lateSubmissionSummary

			if dbSumm.LateGraded != nil {
				l.Graded = *dbSumm.LateGraded
				summ.Graded += l.Graded
			}

			if dbSumm.LatePendingReview != nil {
				l.PendingReview = *dbSumm.LatePendingReview
				summ.PendingReview += l.PendingReview
			}

			if dbSumm.LateSubmitted != nil {
				l.Submitted = *dbSumm.LateSubmitted
				summ.Submitted += l.Submitted
			}

			if dbSumm.LateUnsubmitted != nil {
				l.Unsubmitted = *dbSumm.LateUnsubmitted
				summ.Unsubmitted += l.Unsubmitted
			}

			l.Total = l.Graded +
				l.PendingReview +
				l.Submitted +
				l.Unsubmitted

			summ.Total = summ.Graded +
				summ.PendingReview +
				summ.Submitted +
				summ.Unsubmitted

			if lateCount {
				// late divided by all, multiplied by 100 to get a percentage
				l.Percent = float64(l.Total) / float64(l.Total+summ.Total) * 100
				summ.Late = &l
			}

			resp.SubmissionSummary[uID] = summ
		}
	} else {
		var reqIDs []string
		for _, id := range requestedUserIDs {
			reqIDs = append(reqIDs, strconv.Itoa(int(id)))
		}

		canvasSubs, err := getCanvasCourseSubmissions(rd, &getCanvasCourseSubmissionsRequest{
			courseID:   courseID,
			studentIDs: reqIDs,
		})
		if errors.Is(err, canvasErrorInvalidAccessTokenError) || errors.Is(err, canvasErrorInsufficientScopesOnAccessTokenError) {
			handleError(w, GradesErrorResponse{
				Error:  gradesErrorRevokedToken,
				Action: gradesErrorActionRedirectToOAuth,
			}, http.StatusForbidden)
			return
		} else if errors.Is(err, canvasErrorUnknownError) {
			handleError(w, gradesErrorUnknownCanvasErrorResponse, util.CanvasProxyErrorCode)
			return
		} else if err != nil {
			handleISE(w, fmt.Errorf("error getting canvas course submissions for summary: %w", err))
			return
		}

		go saveSubmissionsToDB(*canvasSubs, uint64(cID))

		for _, sub := range *canvasSubs {
			s := resp.SubmissionSummary[sub.UserID]
			l := s.Late

			if lateCount {
				if l == nil {
					l = &lateSubmissionSummary{}
				}
			}

			switch submissions.WorkflowState(sub.WorkflowState) {
			case submissions.WorkflowStateGraded:
				s.Graded++
				s.Total++
				if sub.Late && lateCount {
					l.Graded++
					l.Total++
				}
			case submissions.WorkflowStatePendingReview:
				s.PendingReview++
				s.Total++
				if sub.Late && lateCount {
					l.PendingReview++
					l.Total++
				}
			case submissions.WorkflowStateSubmitted:
				s.Submitted++
				s.Total++
				if sub.Late && lateCount {
					l.Submitted++
					l.Total++
				}
			case submissions.WorkflowStateUnsubmitted:
				s.Unsubmitted++
				s.Total++
				if sub.Late && lateCount {
					l.Unsubmitted++
					l.Total++
				}
			}

			if lateCount {
				l.Percent = float64(l.Total) / float64(l.Total+s.Total) * 100
				s.Late = l
			}

			resp.SubmissionSummary[sub.UserID] = s
		}
	}

	sendJSON(w, &resp)
	return
}
