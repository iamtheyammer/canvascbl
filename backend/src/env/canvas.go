package env

import (
	"fmt"
	"strconv"
)

var (
	CanvasOAuth2ClientID     = getEnvOrPanic("CANVAS_OAUTH2_CLIENT_ID")
	CanvasOAuth2ClientSecret = getEnvOrPanic("CANVAS_OAUTH2_CLIENT_SECRET")
	CanvasOAuth2Subdomain    = getEnvOrPanic("CANVAS_OAUTH2_SUBDOMAIN")
	CanvasOAuth2RedirectURI  = getEnvOrPanic("CANVAS_OAUTH2_REDIRECT_URI")
	CanvasOAuth2SuccessURI   = getEnvOrPanic("CANVAS_OAUTH2_SUCCESS_URI")

	CanvasCurrentEnrollmentTermID = getCanvasCurrentEnrollmentTermID()
)

func getCanvasCurrentEnrollmentTermID() int {
	etID := getEnv("CANVAS_CURRENT_ENROLLMENT_TERM_ID", "0")

	etIDInt, err := strconv.Atoi(etID)
	if err != nil {
		panic(fmt.Errorf("error converting CANVAS_CURRENT_ENROLLMENT_TERM_ID to an int: %w", err))
	}

	return etIDInt
}
