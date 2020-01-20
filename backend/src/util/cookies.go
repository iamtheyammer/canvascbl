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
		Expires:  time.Now().Add(time.Hour * 312),
	}
	http.SetCookie(w, &c)
}
