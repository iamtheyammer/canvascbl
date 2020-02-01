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

func GetAverageGradeForCourseHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	courseID := ps.ByName("courseID")
	if len(courseID) < 1 {
		util.SendBadRequest(w, "missing courseID as url param")
		return
	}

	if !util.ValidateIntegerString(courseID) {
		util.SendBadRequest(w, "courseID doesn't look like a number")
		return
	}

	cID, err := strconv.Atoi(courseID)
	if err != nil {
		util.SendBadRequest(w, "unable to convert courseID url param to int")
		return
	}

	session := middlewares.Session(w, req, true)
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

	avg, numFactors, err := db.GetMemoizedAverageGradeForCourse(uint64(cID), usersToList)
	if err != nil {
		util.HandleError(errors.Wrap(err, "error getting memoized average for course"))
		util.SendInternalServerError(w)
		return
	}

	ret := struct {
		NumFactors   int    `json:"numFactors"`
		AverageGrade string `json:"averageGrade"`
		Error        string `json:"error"`
	}{}

	if avg == nil || numFactors == nil {
		ret.NumFactors = 0
		ret.AverageGrade = "N/A"
		ret.Error = errorGettingFactorsMessage
	}

	// ensures numFactors hasn't been set
	if len(ret.AverageGrade) < 1 && *numFactors <= minFactorsInAverage {
		ret.Error = notEnoughFactorsMessage
		ret.NumFactors = -1
		ret.AverageGrade = "N/A"
	}

	if len(ret.AverageGrade) < 1 {
		ret.NumFactors = int(*numFactors)
		ret.AverageGrade = util.ConvertGradeAverageToString(*avg)
	}

	jret, err := json.Marshal(ret)
	if err != nil {
		util.HandleError(errors.Wrap(err, "error marshaling average grade for class"))
		util.SendInternalServerError(w)
		return
	}

	util.SendJSONResponse(w, jret)
	return
}
