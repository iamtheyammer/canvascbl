package gradesapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/iamtheyammer/canvascbl/backend/src/oauth2"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type listCoursesResponse struct {
	Courses               *canvasCoursesResponse  `json:"courses"`
	DistanceLearningPairs *[]distanceLearningPair `json:"distance_learning_pairs,omitempty"`
}

type listEnrollmentsResponse struct {
	Enrollments []canvasFullEnrollment `json:"enrollments"`
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

	userID, rdP, sess := authorizer(w, r, []oauth2.Scope{oauth2.ScopeEnrollments, oauth2.ScopeGrades}, &oauth2.AuthorizerAPICall{
		Method:    "GET",
		RoutePath: "courses/:courseID/enrollments",
		Query:     &r.URL.RawQuery,
	})
	if (userID == nil || rdP == nil) && sess == nil {
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

	var es canvasEnrollmentsResponse
	_, err := handleRequestWithTokenRefresh(func(reqD *requestDetails) error {
		tEnrollments, tErr := getCanvasCourseEnrollments(*reqD, courseID, types, states, []string{"avatar_url", "observed_users"})
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
