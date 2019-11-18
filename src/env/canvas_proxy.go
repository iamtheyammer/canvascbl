package env

import "strings"

var (
	ProxyAllowedCORSOrigins = getEnvOrPanic("CANVAS_PROXY_ALLOWED_CORS_ORIGINS")
	ProxyAllowedSubdomains  = strings.Split(getEnvOrPanic("CANVAS_PROXY_ALLOWED_SUBDOMAINS"), ",")
	DefaultSubdomain        = getEnv("CANVAS_PROXY_DEFAULT_SUBDOMAIN", "canvas")
)
