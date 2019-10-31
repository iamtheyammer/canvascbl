package canvasapis

import (
	"github.com/iamtheyammer/canvascbl/backend/src/canvasapis/services/users"
	"github.com/iamtheyammer/canvascbl/backend/src/db"
	"github.com/iamtheyammer/canvascbl/backend/src/email"
	"github.com/iamtheyammer/canvascbl/backend/src/env"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var cookieNoSameSiteRegex = regexp.MustCompile("(^.*iPhone; CPU iPhone OS 1[0-3].*$|^.*iPad; CPU OS 1[0-3].*$|^.*iPod touch; CPU iPhone OS 1[0-3].*$|^.*Macintosh; Intel Mac OS X.*Version\\/1[0-3].*Safari.*$)")

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

	shouldGenerateSession := strings.ToLower(r.URL.Query().Get("generateSession")) == "true"
	if shouldGenerateSession {
		ss, err := db.UpsertProfileAndGenerateSession(&body)
		if err != nil {
			util.SendInternalServerError(w)
			return
		}

		// for API use, possibly (can't hurt!)
		w.Header().Set("X-Session-String", *ss)

		secure := true

		if env.Env == env.EnvironmentDevelopment {
			secure = false
		}

		c := http.Cookie{
			Name:  "session_string",
			Value: *ss,
			Path:  "/",
			// chrome will only deliver cross-site cookies if this is secure in production-- but it won't work in
			// local environments, so we set this depending on the environment (false in dev, otherwise true)
			// see: https://www.chromestatus.com/feature/5088147346030592, https://web.dev/samesite-cookies-explained/
			Secure: secure,
			// chrome will only deliver cross-site cookies if this is set to none
			// 2 weeks
			Expires: time.Now().Add(time.Hour * 312),
		}

		// this ugly workaround is for a Safari 12 bug that causes SameSite=None cookies to be treated as
		// SameSite=Strict cookies. Chrome, however, will only deliver cookies cross-origin if SameSite=None,
		// so this makes it so that incompatible Safari installations don't get the cookie but everyone else
		// does.
		// see: https://bugs.webkit.org/show_bug.cgi?id=198181
		if !cookieNoSameSiteRegex.MatchString(r.Header.Get("User-Agent")) {
			c.SameSite = http.SameSiteNoneMode
		}

		http.SetCookie(w, &c)
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
