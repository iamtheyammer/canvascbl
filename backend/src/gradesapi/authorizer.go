package gradesapi

import (
	"errors"
	"fmt"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/sessions"
	"github.com/iamtheyammer/canvascbl/backend/src/middlewares"
	"github.com/iamtheyammer/canvascbl/backend/src/oauth2"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"net/http"
	"strings"
)

/*
authorizer authorizes an API call.

It returns the user's ID, a requestDetails object and a verified session,
if the call was authorized via session-- DO NOT EXPECT THIS AND CHECK
FOR NIL!
*/
func authorizer(
	w http.ResponseWriter,
	r *http.Request,
	scopes []oauth2.Scope,
	call *oauth2.AuthorizerAPICall,
) (
	*uint64,
	*requestDetails,
	*sessions.VerifiedSession,
	*util.APIErrorContext,
) {
	var (
		at, tokenIsOK = middlewares.Bearer(w, r, false)
		session       *sessions.VerifiedSession
		rd            requestDetails
		userID        uint64
		errCtx        util.APIErrorContext
	)

	if !tokenIsOK {
		handleError(w, GradesErrorResponse{
			Error: gradesErrorInvalidAccessToken,
		}, http.StatusUnauthorized)
		return nil, nil, nil, nil
	}

	// copy headers
	for n, v := range r.Header {
		if errCtx.RequestHeaders == nil {
			errCtx.RequestHeaders = make(map[string][]string, len(r.Header))
		}

		switch strings.ToLower(n) {
		case "cookie":
			errCtx.RequestHeaders[n] = []string{"<REDACTED>"}
		case "authorization":
			errCtx.RequestHeaders[n] = []string{"<REDACTED>"}
		default:
			errCtx.RequestHeaders[n] = v
		}
	}

	if len(at) < 1 {
		// session time
		session = middlewares.Session(w, r, true)
		if session == nil {
			return nil, nil, nil, nil
		}

		errCtx.AuthorizationMethod = "session"
		userID = session.UserID
	} else {
		// oauth2
		grant, err := oauth2.Authorizer(at, scopes, call)
		if err != nil {
			if errors.Is(err, oauth2.GrantMissingScopeError) {
				handleError(w, GradesErrorResponse{
					Error: gradesErrorUnauthorizedScope,
				}, http.StatusUnauthorized)
				return nil, nil, nil, nil
			}

			if errors.Is(err, oauth2.InvalidAccessTokenError) {
				handleError(w, GradesErrorResponse{
					Error: oauth2.InvalidAccessTokenError.Error(),
				}, http.StatusForbidden)
				return nil, nil, nil, nil
			}

			handleISE(w, fmt.Errorf("error using oauth2.Authorizer: %w", err))
			return nil, nil, nil, nil
		}

		errCtx.AuthorizationMethod = "oauth2_bearer"
		errCtx.AddCustomField("oauth2_grant_id", grant.ID)
		userID = grant.UserID
	}

	rd, err := rdFromUserID(userID)
	if err != nil {
		handleISE(w, fmt.Errorf("error getting rd from user id: %w", err))
		return nil, nil, nil, nil
	}

	if rd.TokenID < 1 {
		handleError(w, GradesErrorResponse{
			Error:  gradesErrorNoTokens,
			Action: gradesErrorActionRedirectToOAuth,
		}, http.StatusForbidden)
		return nil, nil, nil, nil
	}

	if call != nil {
		errCtx.Path = call.RoutePath
		errCtx.Method = call.Method
		if call.Query != nil {
			errCtx.Query = *call.Query
		}
	}

	for _, s := range scopes {
		errCtx.Scopes = append(errCtx.Scopes, string(s))
	}

	errCtx.UserID = userID

	return &userID, &rd, session, &errCtx
}
