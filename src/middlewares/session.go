package middlewares

import (
	"github.com/iamtheyammer/canvascbl/backend/src/db"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/sessions"
	"github.com/iamtheyammer/canvascbl/backend/src/env"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

// Middlewares have the same function signature as a request handler, however they return a pointer to an object.
// If that object is nil, the middleware has sent a response and your route handler should return.

// Session ensures that a session string is present and that the user associated with the
// session string has a valid subscription.
func Session(w http.ResponseWriter, req *http.Request) *sessions.VerifiedSession {
	var sessionString string

	cSession, err := req.Cookie("session_string")
	if err != nil {
		if err != http.ErrNoCookie {
			util.SendBadRequest(w, "there was an error parsing your cookie")
			return nil
		}
	} else {
		sessionString = cSession.Value
	}

	hSession := req.Header.Get("X-Session-String")
	if len(hSession) > 0 {
		if len(sessionString) > 0 && hSession != sessionString {
			util.SendBadRequest(w, "you sent a cookie session and a header session and they don't match")
			return nil
		}
		sessionString = hSession
	}

	if len(sessionString) < 1 {
		util.SendUnauthorized(w, "no session string (pass it in via the session_string cookie or the X-Session-String header)")
		return nil
	}

	referer := req.Header.Get("Referer")
	origin := req.Header.Get("Origin")
	if (len(referer) > 0 && !strings.HasPrefix(referer, env.ProxyAllowedCORSOrigins)) ||
		(len(origin) > 0 && origin != env.ProxyAllowedCORSOrigins) {
		util.SendBadRequest(w, "csrf detected-- if you're trying to use this as an api, either remove "+
			"the Origin and Referer headers or set them to the value of the Access-Control-Allow-Origin header "+
			"from this request")
		return nil
	}

	if !util.ValidateUUIDString(sessionString) {
		util.SendBadRequest(w, "the session string sent does not look like a valid session string")
		return nil
	}

	sessionInfo, err := db.VerifySession(sessionString)
	if err != nil {
		util.HandleError(errors.Wrap(err, "error verifying session"))
		util.SendUnauthorized(w, "unable to verify your session")
		return nil
	}

	if sessionInfo == nil {
		util.SendUnauthorized(w, "invalid session string")
		return nil
	}

	if sessionInfo.SessionIsExpired {
		util.SendUnauthorized(w, "expired session")
		return nil
	}

	return sessionInfo
}
