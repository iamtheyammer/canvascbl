package util

import (
	"errors"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/iamtheyammer/canvascbl/backend/src/env"
	"github.com/lib/pq"
)

func HandleError(err error) {
	e := err

	var pqErr *pq.Error
	if errors.As(e, &pqErr) {
		e = fmt.Errorf("postgres extra error data: detail: %s, hint: %s; %w", pqErr.Detail, pqErr.Hint, e)
	}

	fmt.Println(e.Error())

	if env.Env != env.EnvironmentDevelopment {
		sentry.WithScope(func(scope *sentry.Scope) {
			if pqErr != nil {
				scope.SetTag("from", "postgres")

				if len(pqErr.Detail) > 0 {
					scope.SetContext("postgres_detail", pqErr.Detail)
				}

				if len(pqErr.Hint) > 0 {
					scope.SetContext("postgres_hint", pqErr.Hint)
				}
			}

			// using err here because we already captured postgres data above
			sentry.CaptureException(err)
		})
	}
}
