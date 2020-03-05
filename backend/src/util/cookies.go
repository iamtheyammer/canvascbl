package util

import (
	"github.com/iamtheyammer/canvascbl/backend/src/env"
	"net/http"
	"time"
)

// AddSessionToResponse adds the newly generated session string to a response.
func AddSessionToResponse(w http.ResponseWriter, ss string) {
	secure := true
	sameSite := http.SameSiteStrictMode

	if env.Env == env.EnvironmentDevelopment {
		secure = false
		sameSite = http.SameSiteNoneMode
	}

	c := http.Cookie{
		Name:     "session_string",
		Value:    ss,
		Path:     "/",
		SameSite: sameSite,
		Secure:   secure,
		// no longer expiring because we can use "expired session" to intent=reauth instead of auth.
		//Expires:  time.Now().Add(time.Hour * 312),
		// set expires to a very long time to not have them be session cookies
		Expires: time.Now().Add(time.Hour * 5000),
	}
	http.SetCookie(w, &c)
}
