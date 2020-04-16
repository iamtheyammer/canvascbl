package env

import (
	"fmt"
	"strings"
)

var (
	// HTTPPort is the port for the HTTP Server to run on.
	HTTPPort = fmt.Sprintf(":%s", getEnv("PORT", "8000"))
	// BaseURL is the API's URL, like "https://api.canvascbl.com"
	BaseURL = getEnvOrPanic("BASE_URL")
	// ScriptKey is used for autonomous tasks that need a key to run.
	ScriptKey = getEnv("SCRIPT_KEY", "default_script_key")
	// Allowed CORS Origins for Access-Control-Allow-Origin. Comma-separated.
	ProxyAllowedCORSOrigins = getProxyAllowedCORSOrgins()
)

func getProxyAllowedCORSOrgins() map[string]struct{} {
	origins := strings.Split(getEnvOrPanic("CANVAS_PROXY_ALLOWED_CORS_ORIGINS"), ",")

	allowed := map[string]struct{}{}
	for _, origin := range origins {
		allowed[origin] = struct{}{}
	}

	return allowed
}
