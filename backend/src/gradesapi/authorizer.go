package gradesapi

import (
	"errors"
	"fmt"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/sessions"
	"github.com/iamtheyammer/canvascbl/backend/src/middlewares"
	"github.com/iamtheyammer/canvascbl/backend/src/oauth2"
	"net/http"
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
) {
	var (
		at, tokenIsOK = middlewares.Bearer(w, r, false)
		session       *sessions.VerifiedSession
		rd            requestDetails
		userID        uint64
	)

	if !tokenIsOK {
		handleError(w, gradesErrorResponse{
			Error: gradesErrorInvalidAccessToken,
		}, http.StatusUnauthorized)
		return nil, nil, nil
	}

	if len(at) < 1 {
		// session time
		session = middlewares.Session(w, r, true)
		if session == nil {
			return nil, nil, nil
		}

		userID = session.UserID
	} else {
		// oauth2
		grant, err := oauth2.Authorizer(at, scopes, call)
		if err != nil {
			if errors.Is(err, oauth2.GrantMissingScopeError) {
				handleError(w, gradesErrorResponse{
					Error: gradesErrorUnauthorizedScope,
				}, http.StatusUnauthorized)
				return nil, nil, nil
			}

			if errors.Is(err, oauth2.InvalidAccessTokenError) {
				handleError(w, gradesErrorResponse{
					Error: oauth2.InvalidAccessTokenError.Error(),
				}, http.StatusForbidden)
				return nil, nil, nil
			}

			handleISE(w, fmt.Errorf("error using oauth2.Authorizer: %w", err))
			return nil, nil, nil
		}

		userID = grant.UserID
	}

	rd, err := rdFromUserID(userID)
	if err != nil {
		handleISE(w, fmt.Errorf("error getting rd from user id: %w", err))
		return nil, nil, nil
	}

	if rd.TokenID < 1 {
		handleError(w, gradesErrorResponse{
			Error:  gradesErrorNoTokens,
			Action: gradesErrorActionRedirectToOAuth,
		}, http.StatusForbidden)
		return nil, nil, nil
	}

	return &userID, &rd, session
}
