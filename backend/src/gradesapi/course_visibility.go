package gradesapi

import (
	"encoding/json"
	"fmt"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/courses"
	"github.com/iamtheyammer/canvascbl/backend/src/oauth2"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

type courseVisibilityResponse struct {
	Hidden bool `json:"canvascbl_hidden"`
}

// HideCourseHandler hides courses.
func HideCourseHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	courseID := ps.ByName("courseID")
	if len(courseID) < 1 || !util.ValidateIntegerString(courseID) {
		util.SendBadRequest(w, "invalid or missing course_id as url param")
		return
	}

	cID, err := strconv.Atoi(courseID)
	if err != nil {
		handleISE(w, fmt.Errorf("error converting course ID to an int in HideCourseHandler: %w", err))
		return
	}

	userID, rdP, sess, errCtx := authorizer(w, r, []oauth2.Scope{oauth2.ScopeCourses}, &oauth2.AuthorizerAPICall{
		Method:    "PUT",
		RoutePath: "courses/:courseID/hide",
	})
	if (userID == nil || rdP == nil || errCtx == nil) && sess == nil {
		return
	}

	errCtx.AddCustomField("course_id", courseID)

	err = courses.Hide(db, &courses.HideRequest{
		UserID:   *userID,
		CourseID: uint64(cID),
	})
	if err != nil {
		handleISE(w, errCtx.Apply(fmt.Errorf("error hiding a course: %w", err)))
		return
	}

	j, err := json.Marshal(&courseVisibilityResponse{Hidden: true})
	if err != nil {
		handleISE(w, errCtx.Apply(fmt.Errorf("error marshaling hide course json: %w", err)))
		return
	}

	util.SendJSONResponse(w, j)
}

// ShowCourseHandler shows courses.
func ShowCourseHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	courseID := ps.ByName("courseID")
	if len(courseID) < 1 || !util.ValidateIntegerString(courseID) {
		util.SendBadRequest(w, "invalid or missing course_id as url param")
		return
	}

	cID, err := strconv.Atoi(courseID)
	if err != nil {
		handleISE(w, fmt.Errorf("error converting course ID to an int in HideCourseHandler: %w", err))
		return
	}

	userID, rdP, sess, errCtx := authorizer(w, r, []oauth2.Scope{oauth2.ScopeCourses}, &oauth2.AuthorizerAPICall{
		Method:    "DELETE",
		RoutePath: "courses/:courseID/hide",
	})
	if (userID == nil || rdP == nil || errCtx == nil) && sess == nil {
		return
	}

	errCtx.AddCustomField("course_id", courseID)

	err = courses.Show(db, &courses.HideRequest{
		UserID:   *userID,
		CourseID: uint64(cID),
	})
	if err != nil {
		handleISE(w, errCtx.Apply(fmt.Errorf("error showing a course: %w", err)))
		return
	}

	j, err := json.Marshal(&courseVisibilityResponse{Hidden: false})
	if err != nil {
		handleISE(w, errCtx.Apply(fmt.Errorf("error marshaling show course json: %w", err)))
		return
	}

	util.SendJSONResponse(w, j)
}
