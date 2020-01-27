package plus

import (
	"encoding/json"
	"github.com/iamtheyammer/canvascbl/backend/src/env"
	"github.com/iamtheyammer/canvascbl/backend/src/middlewares"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

type sessionInformation struct {
	UserID               uint64 `json:"userId"`
	CanvasUserID         uint64 `json:"canvasUserId"`
	Status               int    `json:"status"`
	Email                string `json:"email"`
	HasValidSubscription bool   `json:"hasValidSubscription"`
	SubscriptionStatus   string `json:"subscriptionStatus"`
}

func GetSessionInformationHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	sess := middlewares.Session(w, req)
	if sess == nil {
		return
	}

	sessj, err := json.Marshal(sessionInformation{
		UserID:               sess.UserID,
		CanvasUserID:         sess.CanvasUserID,
		Status:               sess.UserStatus,
		Email:                sess.Email,
		HasValidSubscription: sess.HasValidSubscription,
		SubscriptionStatus:   sess.SubscriptionStatus,
	})
	if err != nil {
		util.HandleError(errors.Wrap(err, "error marshaling get session information struct"))
		return
	}

	util.SendJSONResponse(w, sessj)
	return
}

func ClearSessionHandler(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	secure := true
	sameSite := http.SameSiteStrictMode

	if env.Env == env.EnvironmentDevelopment {
		secure = false
		sameSite = http.SameSiteNoneMode
	}

	c := http.Cookie{
		Name:     "session_string",
		Value:    "",
		Path:     "/",
		SameSite: sameSite,
		Secure:   secure,
		Expires:  time.Now().Add(-time.Hour),
	}
	http.SetCookie(w, &c)
}
