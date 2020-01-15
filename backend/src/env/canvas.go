package env

var (
	CanvasOAuth2ClientID     = getEnvOrPanic("CANVAS_OAUTH2_CLIENT_ID")
	CanvasOAuth2ClientSecret = getEnvOrPanic("CANVAS_OAUTH2_CLIENT_SECRET")
	CanvasOAuth2Subdomain    = getEnvOrPanic("CANVAS_OAUTH2_SUBDOMAIN")
	CanvasOAuth2RedirectURI  = getEnvOrPanic("CANVAS_OAUTH2_REDIRECT_URI")
	CanvasOAuth2SuccessURI   = getEnvOrPanic("CANVAS_OAUTH2_SUCCESS_URI")
)
