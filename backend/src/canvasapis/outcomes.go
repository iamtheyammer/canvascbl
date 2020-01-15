package canvasapis

import (
	"github.com/iamtheyammer/canvascbl/backend/src/canvasapis/services/outcomes"
	"github.com/iamtheyammer/canvascbl/backend/src/db"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

/*
GetOutcomeByIDHandler handles getting outcomes by ID.
*/
func GetOutcomeByIDHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	outcomeID := ps.ByName("outcomeID")
	if len(outcomeID) < 1 {
		util.SendBadRequest(w, "missing outcome id as param id")
		return
	}

	if !util.ValidateIntegerString(outcomeID) {
		util.SendBadRequest(w, "invalid outcomeID")
		return
	}

	ok, rd := util.GetRequestDetailsFromRequest(r)

	if !ok {
		util.SendUnauthorized(w, util.RequestDetailsFailedValidationMessage)
		return
	}

	resp, body, err := outcomes.GetByID(rd, outcomeID)
	if err != nil {
		util.SendInternalServerError(w)
		return
	}

	util.HandleCanvasResponse(w, resp, body)

	if resp.StatusCode != http.StatusOK {
		return
	}

	// db
	go db.InsertOutcome(&body)
	return
}
