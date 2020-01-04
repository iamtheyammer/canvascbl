package canvasapis

import (
	"github.com/iamtheyammer/canvascbl/backend/src/canvasapis/services/users"
	"github.com/iamtheyammer/canvascbl/backend/src/db"
	"github.com/iamtheyammer/canvascbl/backend/src/email"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"strings"
)

func GetOwnUserProfileHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ok, rd := util.GetRequestDetailsFromRequest(r)
	if !ok {
		util.SendUnauthorized(w, util.RequestDetailsFailedValidationMessage)
		return
	}

	resp, body, err := users.GetSelfProfile(rd)
	if err != nil {
		util.SendInternalServerError(w)
		return
	}

	if resp.StatusCode != http.StatusOK {
		util.HandleCanvasResponse(w, resp, body)
		return
	}

	shouldGenerateSession := strings.ToLower(r.URL.Query().Get("generateSession")) == "true"
	if shouldGenerateSession {
		ss, err := db.UpsertProfileAndGenerateSession(&body)
		if err != nil {
			util.SendInternalServerError(w)
			return
		}

		util.AddSessionToResponse(w, *ss)
	}

	util.HandleCanvasResponse(w, resp, body)

	if resp.StatusCode != http.StatusOK {
		return
	}

	// db
	go email.SendWelcomeIfNecessary(&body)

	if !shouldGenerateSession {
		go db.UpsertProfile(&body)
	}

	return
}

func GetOwnObserveesHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userID := r.URL.Query().Get("user_id")
	if len(userID) < 1 {
		util.SendBadRequest(w, "missing user_id as query param")
		return
	}

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		util.HandleError(errors.Wrap(err, "error turning user_id into an int"))
		util.SendBadRequest(w, "invalid user_id as query param")
		return
	}

	ok, rd := util.GetRequestDetailsFromRequest(r)
	if !ok {
		util.SendUnauthorized(w, util.RequestDetailsFailedValidationMessage)
		return
	}

	resp, body, err := users.GetUserObservees(rd, userID)
	if err != nil {
		util.SendInternalServerError(w)
		return
	}

	util.HandleCanvasResponse(w, resp, body)

	if resp.StatusCode == http.StatusOK {
		go db.HandleObservees(&body, uint64(userIDInt))
	}
	return
}
