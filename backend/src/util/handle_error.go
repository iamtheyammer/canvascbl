package util

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/iamtheyammer/canvascbl/backend/src/env"
)

func HandleError(err error) {
	fmt.Println(err.Error())

	if env.Env != env.EnvironmentDevelopment {
		sentry.CaptureException(err)
	}
}
