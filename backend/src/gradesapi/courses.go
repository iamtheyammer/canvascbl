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
	if errors.Is(err, canvasErrorInvalidAccessTokenError) {
		handleError(w, GradesErrorResponse{
			Error:  gradesErrorRevokedToken,
			Action: gradesErrorActionRedirectToOAuth,
		}, http.StatusForbidden)
		return
	} else if errors.Is(err, canvasErrorInvalidAccessTokenError) {
		handleError(w, GradesErrorResponse{
			Error:  gradesErrorRefreshedTokenError,
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
