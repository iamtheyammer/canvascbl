package middlewares

import (
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/sessions"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"net/http"
)

func IsAdmin(w http.ResponseWriter, req *http.Request, session *sessions.VerifiedSession) bool {
	if session.UserStatus != 2 {
		util.SendUnauthorized(w, "you aren't an admin")
		return true
	}

	return false
}
