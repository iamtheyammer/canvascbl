package env

var (
	OAuth2ConsentURL = getEnvOrPanic("OAUTH2_CONSENT_URL")
)
