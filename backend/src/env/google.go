package env

import "strings"

var (
	GoogleOAuth2ClientID     = getEnvOrPanic("GOOGLE_OAUTH2_CLIENT_ID")
	GoogleOAuth2ClientSecret = getEnvOrPanic("GOOGLE_OAUTH2_CLIENT_SECRET")
	GoogleOAuth2RedirectURI  = getEnvOrPanic("GOOGLE_OAUTH2_REDIRECT_URI")
	GoogleOAuth2AllowedHDs   = strings.Split(getEnvOrPanic("GOOGLE_OAUTH2_ALLOWED_HDS"), ",")
)
