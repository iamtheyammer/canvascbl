package canvasapis

import (
	"github.com/iamtheyammer/canvascbl/backend/src/canvasapis/services/courses"
	"github.com/iamtheyammer/canvascbl/backend/src/db"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func GetCoursesHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ok, rd := util.GetRequestDetailsFromRequest(r)
	if !ok {
		util.SendUnauthorized(w, util.RequestDetailsFailedValidationMessage)
		return
	}

	resp, body, err := courses.Get(rd)
	if err != nil {
		util.SendInternalServerError(w)
		return
	}

	util.HandleCanvasResponse(w, resp, body)

	if resp.StatusCode != http.StatusOK {
		return
	}

	// db
	go db.UpsertMultipleCourses(&body)

	return
}

func GetOutcomesByCourseHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	courseID := ps.ByName("courseID")
	if len(courseID) < 1 {
		util.SendBadRequest(w, "missing courseID as url param")
		return
	}

	if !util.ValidateIntegerString(courseID) {
		util.SendBadRequest(w, "invalid courseID")
		return
	}

	ok, rd := util.GetRequestDetailsFromRequest(r)
	if !ok {
		util.SendUnauthorized(w, util.RequestDetailsFailedValidationMessage)
		return
	}

	resp, body, err := courses.GetOutcomesByCourse(rd, courseID)
	if err != nil {
		util.SendInternalServerError(w)
		log.Fatal(err)
		return
	}

	util.HandleCanvasResponse(w, resp, body)
	return
}

func GetOutcomesByCourseAndOutcomeGroupHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	courseID := ps.ByName("courseID")
	if len(courseID) < 1 {
		util.SendBadRequest(w, "missing courseID as url param")
		return
	}

	if !util.ValidateIntegerString(courseID) {
		util.SendBadRequest(w, "invalid courseID")
		return
	}

	outcomeGroupID := ps.ByName("outcomeGroupID")
	if len(outcomeGroupID) < 1 {
		util.SendBadRequest(w, "missing outcomeGroupID as url param")
		return
	}

	if !util.ValidateIntegerString(outcomeGroupID) {
		util.SendBadRequest(w, "invalid outcomeGroupID")
		return
	}

	ok, rd := util.GetRequestDetailsFromRequest(r)
	if !ok {
		util.SendUnauthorized(w, util.RequestDetailsFailedValidationMessage)
		return
	}

	resp, body, err := courses.GetOutcomesByCourseAndOutcomeGroup(rd, courseID, outcomeGroupID)
	if err != nil {
		util.SendInternalServerError(w)
		return
	}

	util.HandleCanvasResponse(w, resp, body)
	return
}

func GetOutcomeResultsByCourseHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	courseID := ps.ByName("courseID")
	if len(courseID) < 1 {
		util.SendBadRequest(w, "missing courseID as url param")
		return
	}

	if !util.ValidateIntegerString(courseID) {
		util.SendBadRequest(w, "invalid courseID")
		return
	}

	userIDs := r.URL.Query()["user_ids[]"]
	if len(userIDs) < 1 {
		util.SendBadRequest(w, "missing user_ids[] as param user_ids[]")
		return
	}

	for _, uID := range userIDs {
		if !util.ValidateIntegerString(uID) {
			util.SendBadRequest(w, "one of your user_ids[] is invalid")
			return
		}
	}

	ok, rd := util.GetRequestDetailsFromRequest(r)
	if !ok {
		util.SendUnauthorized(w, util.RequestDetailsFailedValidationMessage)
		return
	}

	includes := r.URL.Query().Get("include[]")

	if len(includes) > 1 {
		if !util.ValidateIncludes(includes) {
			util.SendBadRequest(w, "invalid includes")
			return
		}
	}

	resp, body, err := courses.GetOutcomeResultsByCourse(rd, courseID, userIDs, includes)
	if err != nil {
		util.SendInternalServerError(w)
		return
	}

	util.HandleCanvasResponse(w, resp, body)

	if resp.StatusCode != http.StatusOK {
		return
	}

	// db
	go db.InsertMultipleOutcomeResults(&body, &courseID)

	return
}

func GetOutcomeRollupsByCourseHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	courseID := ps.ByName("courseID")
	if len(courseID) < 1 {
		util.SendBadRequest(w, "missing courseID as url param")
		return
	}

	if !util.ValidateIntegerString(courseID) {
		util.SendBadRequest(w, "invalid courseID")
		return
	}

	userIDs := r.URL.Query()["user_ids[]"]
	if len(userIDs) < 1 {
		util.SendBadRequest(w, "missing user_ids[] as param user_ids[]")
		return
	}

	for _, uID := range userIDs {
		if !util.ValidateIntegerString(uID) {
			util.SendBadRequest(w, "one of your user_ids[] is invalid")
			return
		}
	}

	ok, rd := util.GetRequestDetailsFromRequest(r)
	if !ok {
		util.SendUnauthorized(w, util.RequestDetailsFailedValidationMessage)
		return
	}

	includes := r.URL.Query().Get("include[]")

	if len(includes) > 1 {
		if !util.ValidateIncludes(includes) {
			util.SendBadRequest(w, "invalid includes")
			return
		}
	}

	resp, body, err := courses.GetOutcomeRollupsByCourse(rd, courseID, userIDs, includes)
	if err != nil {
		util.SendInternalServerError(w)
		return
	}

	util.HandleCanvasResponse(w, resp, body)

	if resp.StatusCode != http.StatusOK {
		return
	}

	// send to db; using go funcs so they run at the same time

	// grades
	go db.InsertGrade(&body, &courseID, &userIDs)
	// rollups
	go db.InsertMultipleOutcomeRollups(&body, &courseID)

	return
}

func GetAssignmentsByCourseHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	courseID := ps.ByName("courseID")
	if len(courseID) < 1 {
		util.SendBadRequest(w, "missing courseID as url param")
		return
	}

	if !util.ValidateIntegerString(courseID) {
		util.SendBadRequest(w, "invalid courseID")
		return
	}

	ok, rd := util.GetRequestDetailsFromRequest(r)
	if !ok {
		util.SendUnauthorized(w, util.RequestDetailsFailedValidationMessage)
		return
	}

	includes := r.URL.Query().Get("include[]")

	if len(includes) > 1 {
		if !util.ValidateIncludes(includes) {
			util.SendBadRequest(w, "invalid includes")
			return
		}
	}

	resp, body, err := courses.GetAssignmentsByCourse(rd, courseID, includes)
	if err != nil {
		util.SendInternalServerError(w)
		return
	}

	util.HandleCanvasResponse(w, resp, body)

	if resp.StatusCode != http.StatusOK {
		return
	}

	// db
	go db.InsertMultipleAssignments(&body)
	return
}

func GetOutcomeAlignmentsByCourseHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	courseID := ps.ByName("courseID")
	if len(courseID) < 1 {
		util.SendBadRequest(w, "missing courseID as url param")
		return
	}

	if !util.ValidateIntegerString(courseID) {
		util.SendBadRequest(w, "invalid courseID")
		return
	}

	userID := r.URL.Query().Get("userId")
	if len(userID) < 1 {
		util.SendBadRequest(w, "missing userId as param userId")
		return
	}

	if !util.ValidateIntegerString(userID) {
		util.SendBadRequest(w, "invalid userId")
		return
	}

	ok, rd := util.GetRequestDetailsFromRequest(r)
	if !ok {
		util.SendUnauthorized(w, util.RequestDetailsFailedValidationMessage)
		return
	}

	resp, body, err := courses.GetOutcomeAlignmentsByCourse(rd, courseID, userID)
	if err != nil {
		util.SendInternalServerError(w)
		return
	}

	util.HandleCanvasResponse(w, resp, body)
}
