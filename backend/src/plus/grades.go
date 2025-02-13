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
	"time"
)

type userGrade struct {
	ID           uint64 `json:"id"`
	CanvasUserID uint64 `json:"canvasUserId"`
	CourseID     uint64 `json:"courseId"`
	Grade        string `json:"grade"`
	InsertedAt   int64  `json:"insertedAt"`
}

func GetPreviousGradesHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	session := middlewares.Session(w, req, true)
	if session == nil {
		return
	}
	if hvs := middlewares.ValidSubscription(w, req, session); hvs {
		return
	}

	obsP, err := db.ListObservees(&users.ListObserveesRequest{
		ObserverCanvasUserID: session.CanvasUserID,
		ActiveOnly:           true,
	})
	if err != nil {
		util.HandleError(errors.Wrap(err, "error getting user observees"))
		util.SendInternalServerError(w)
		return
	}

	obs := *obsP

	var userIdsToGetGradesFor []uint64

	if len(obs) > 1 {
		for _, o := range obs {
			userIdsToGetGradesFor = append(userIdsToGetGradesFor, o.CanvasUserID)
		}
	} else {
		userIdsToGetGradesFor = append(userIdsToGetGradesFor, session.CanvasUserID)
	}

	gsP, err := db.GetGradesForUserBeforeDate(
		userIdsToGetGradesFor,
		time.Now().Add(-(time.Minute * 5)),
	)

	if err != nil {
		util.HandleError(errors.Wrap(err, "error getting grades for user before date"))
		util.SendInternalServerError(w)
		return
	}

	gs := *gsP

	if len(gs) < 1 {
		util.SendJSONResponse(w, []byte("[]"))
		return
	}

	var ugs []userGrade

	for _, g := range gs {
		ug := userGrade{
			ID:           g.ID,
			CanvasUserID: g.UserCanvasID,
			CourseID:     g.CourseID,
			Grade:        g.Grade,
			InsertedAt:   g.InsertedAt.Unix(),
		}
		ugs = append(ugs, ug)
	}

	jugs, err := json.Marshal(ugs)
	if err != nil {
		util.HandleError(errors.Wrap(err, "error marshaling user grades into json"))
		util.SendInternalServerError(w)
		return
	}

	util.SendJSONResponse(w, jugs)
	return
}
