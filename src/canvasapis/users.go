package canvasapis

import (
	"github.com/iamtheyammer/canvascbl/backend/src/canvasapis/services/users"
	"github.com/iamtheyammer/canvascbl/backend/src/db"
	"github.com/iamtheyammer/canvascbl/backend/src/email"
	"github.com/iamtheyammer/canvascbl/backend/src/env"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"
	"time"
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
		sameSite := http.SameSiteStrictMode

		if env.Env == env.EnvironmentDevelopment {
			secure = false
			sameSite = http.SameSiteNoneMode
		}

		c := http.Cookie{
			Name:     "session_string",
			Value:    *ss,
			Path:     "/",
			SameSite: sameSite,
			Secure:   secure,
			Expires:  time.Now().Add(time.Hour * 312),
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
