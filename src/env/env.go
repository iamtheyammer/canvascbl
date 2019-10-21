package env

import (
	"fmt"
	"os"
	"strings"
)

var HTTPPort = fmt.Sprintf(":%s", getEnv("PORT", "8000"))

var OAuth2ClientID = getEnvOrPanic("CANVAS_OAUTH2_CLIENT_ID")
var OAuth2ClientSecret = getEnvOrPanic("CANVAS_OAUTH2_CLIENT_SECRET")
var OAuth2Subdomain = getEnvOrPanic("CANVAS_OAUTH2_SUBDOMAIN")
var OAuth2RedirectURI = getEnvOrPanic("CANVAS_OAUTH2_REDIRECT_URI")
var OAuth2SuccessURI = getEnvOrPanic("CANVAS_OAUTH2_SUCCESS_URI")

var ProxyAllowedCORSOrigins = getEnvOrPanic("CANVAS_PROXY_ALLOWED_CORS_ORIGINS")
var ProxyAllowedSubdomains = strings.Split(getEnvOrPanic("CANVAS_PROXY_ALLOWED_SUBDOMAINS"), ",")
var DefaultSubdomain = getEnv("CANVAS_PROXY_DEFAULT_SUBDOMAIN", "canvas")

var DatabaseDSN = getEnvOrPanic("DATABASE_DSN")

var StripeAPIKey = getEnvOrPanic("STRIPE_API_KEY")
var StripeWebhookSecret = getEnvOrPanic("STRIPE_WEBHOOK_SECRET")

func getEnvOrPanic(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Sprintf("Missing required environment variable '%v'\n", key))
	}
	return value
}

func getEnv(key string, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}

	return fallback
}
