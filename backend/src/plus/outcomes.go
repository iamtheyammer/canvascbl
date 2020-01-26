package plus

import (
	"encoding/json"
	"github.com/iamtheyammer/canvascbl/backend/src/db"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/users"
	"github.com/iamtheyammer/canvascbl/backend/src/middlewares"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

func GetAverageOutcomeScoreHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	outcomeID := ps.ByName("outcomeID")
	if len(outcomeID) < 1 {
		util.SendBadRequest(w, "missing outcomeID as url param")
		return
	}

	if !util.ValidateIntegerString(outcomeID) {
		util.SendBadRequest(w, "invalid outcomeID as url param")
		return
	}

	oID, err := strconv.Atoi(outcomeID)
	if err != nil {
		util.HandleError(errors.Wrap(err, "error converting outcomeID to int"))
		util.SendInternalServerError(w)
		return
	}

	session := middlewares.Session(w, req)
	if session == nil {
		return
	}

	if hvs := middlewares.ValidSubscription(w, req, session); hvs {
		return
	}

	obsP, err := db.ListObservees(&users.ListObserveesRequest{ObserverCanvasUserID: session.CanvasUserID})
	if err != nil {
		util.HandleError(errors.Wrap(err, "error listing observees"))
		util.SendInternalServerError(w)
		return
	}
	obs := *obsP

	var usersToList []uint64
	if len(obs) > 0 {
		for _, o := range obs {
			usersToList = append(usersToList, o.CanvasUserID)
		}
	} else {
		usersToList = append(usersToList, session.CanvasUserID)
	}

	score, err := db.GetUserMostRecentOutcomeRollupScore(usersToList)
	if err != nil {
		util.HandleError(errors.Wrap(err, "error getting the most recent outcome score for a user"))
		util.SendInternalServerError(w)
		return
	}

	if score == nil {
		util.SendUnauthorized(w, "you don't have a score for that outcome")
		return
	}

	avg, err := db.GetMemoizedOutcomeAverage(uint64(oID))
	if err != nil {
		util.HandleError(errors.Wrap(err, "error getting outcome average"))
		util.SendInternalServerError(w)
		return
	}

	res := struct {
		AverageScore float64 `json:"averageScore"`
		NumFactors   int     `json:"numFactors"`
		Error        string  `json:"error"`
	}{}

	if avg == nil {
		res.AverageScore = -1
		res.NumFactors = -1
		res.Error = errorGettingFactorsMessage
	}

	if avg.NumFactors < minFactorsInAverage && len(res.Error) < 1 {
		res.AverageScore = -1
		res.NumFactors = -1
		res.Error = notEnoughFactorsMessage
	}

	if avg.AverageScore == -1 {
		res.AverageScore = -1
		res.NumFactors = -1
		res.Error = errorGettingAverageMessage
	}

	if len(res.Error) < 1 {
		res.AverageScore = avg.AverageScore
		res.NumFactors = int(avg.NumFactors)
	}

	jret, err := json.Marshal(res)
	if err != nil {
		util.HandleError(errors.Wrap(err, "error marshaling average outcome score json"))
		return
	}

	util.SendJSONResponse(w, jret)
	return
}
