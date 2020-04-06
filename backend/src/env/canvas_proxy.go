package env

var (
	ProxyAllowedCORSOrigins = getEnvOrPanic("CANVAS_PROXY_ALLOWED_CORS_ORIGINS")
)
