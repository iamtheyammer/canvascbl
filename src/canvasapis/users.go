package canvasapis

import (
	"github.com/iamtheyammer/canvascbl/backend/src/canvasapis/services/users"
	"github.com/iamtheyammer/canvascbl/backend/src/db"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/julienschmidt/httprouter"
	"net/http"
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
	}

	util.HandleCanvasResponse(w, resp, body)

	// db
	db.UpsertProfile(&body)

	return
}
