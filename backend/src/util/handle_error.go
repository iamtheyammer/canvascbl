package util

import (
	"errors"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/iamtheyammer/canvascbl/backend/src/env"
	"github.com/lib/pq"
)

// APIErrorContext provides context for Internal API Errors.
type APIErrorContext struct {
	Err                 error
	UserID              uint64
	Path                string
	Method              string
	Query               string
	AuthorizationMethod string
	RequestHeaders      map[string][]string
	Scopes              []string
	customFields        map[string]interface{}
}

// Error returns an error string for compatibility with the builtin error interface.
func (err APIErrorContext) Error() string {
	return err.Err.Error()
}

// AddCustomField allows you to add a custom field to the target.
//
// If two fields with the same name are added, it will preserve both, not modify the original.
// Note that it's expensive to add duplicate fields, so please avoid it.
//
// The idiomatic case for field names is snake_case.
func (err *APIErrorContext) AddCustomField(k string, v interface{}) {
	if err.customFields == nil {
		err.customFields = map[string]interface{}{}
	}

	// if this key already exists, we will try key_2, then key_3, ...
	if _, ok := err.customFields[k]; ok {
		// max 1000, this can't go forever.
		for id := 2; id <= 1000; id++ {
			newKey := fmt.Sprintf("%s_%d", k, id)

			// if this key does NOT already exist
			if _, ok := err.customFields[newKey]; !ok {
				err.customFields[newKey] = v
				break
			}
		}
	} else {
		err.customFields[k] = v
	}
}

// AddCustomFields calls AddCustomField for every field in the provided map.
func (err *APIErrorContext) AddCustomFields(fields map[string]interface{}) {
	for k, v := range fields {
		err.AddCustomField(k, v)
	}
}

// AddUserDetails adds user info to the error context.
func (err *APIErrorContext) AddUserDetails(userID *uint64) {
	if userID != nil {
		err.UserID = *userID
	}
}

// Apply applies the context to the given error.
func (err *APIErrorContext) Apply(e error) error {
	err.Err = fmt.Errorf("%w", e)
	return fmt.Errorf("%w", err)
}

func HandleError(err error) {
	e := err
	unwrappable := err

	var errContext *APIErrorContext
	if errors.As(err, &errContext) {
		e = fmt.Errorf("api context: "+
			"user id: %d, "+
			"path: %s, "+
			"method: %s, "+
			"query: %s, "+
			"authorization method: %s, "+
			"custom fields: %+v; %w",
			errContext.UserID,
			errContext.Path,
			errContext.Method,
			errContext.Query,
			errContext.AuthorizationMethod,
			errContext.customFields,
			errContext.Err,
		)
		unwrappable = errContext.Err
	}

	var pqErr *pq.Error
	if errors.As(unwrappable, &pqErr) {
		e = fmt.Errorf("postgres extra error data: detail: %s, hint: %s; %w", pqErr.Detail, pqErr.Hint, e)
	}

	fmt.Println(e.Error())

	if env.Env != env.EnvironmentDevelopment {
		sentry.WithScope(func(scope *sentry.Scope) {
			if pqErr != nil {
				scope.SetTag("from", "postgres")

				scope.SetContext("Postgres Error", pqErr)
			}

			if errContext != nil {
				scope.SetTag("has_api_error_context", "true")
				scope.SetTag("api_authorization_method", errContext.AuthorizationMethod)

				scope.SetUser(sentry.User{
					ID: fmt.Sprintf("%d", errContext.UserID),
				})

				if errContext.customFields != nil {
					scope.SetContext("API Error Context Custom Fields", errContext.customFields)
				}

				scope.SetContext("API Error Context", errContext)

				// we capture this error so we get a full call stack
				// all context has been added to sentry
				sentry.CaptureException(errContext.Err)
			} else {
				scope.SetTag("has_api_error_context", "false")
				sentry.CaptureException(err)
			}
		})
	}
}
