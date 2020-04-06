package env

import (
	"fmt"
	"strconv"
)

var (
	// CanvasDomain is the domain that Canvas is running on.
	CanvasDomain             = getEnvOrPanic("CANVAS_DOMAIN")
	CanvasOAuth2ClientID     = getEnvOrPanic("CANVAS_OAUTH2_CLIENT_ID")
	CanvasOAuth2ClientSecret = getEnvOrPanic("CANVAS_OAUTH2_CLIENT_SECRET")
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
