package middlewares

import (
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/sessions"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"net/http"
)

func ValidSubscription(w http.ResponseWriter, req *http.Request, session *sessions.VerifiedSession) bool {
	if session.HasValidSubscription == false {
		util.SendUnauthorized(w, "you must have a valid subscription to use this feature")
		return true
	}

	return false
}
