package oauth2

import (
	"errors"
	"fmt"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/oauth2"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"time"
)

type Grant struct {
	ID                 uint64     `json:"id"`
	UserID             uint64     `json:"user_id"`
	Purpose            string     `json:"purpose,omitempty"`
	OAuth2CredentialID uint64     `json:"oauth2_credential_id"`
	RedirectURIID      uint64     `json:"-"`
	AccessToken        string     `json:"-"`
	RefreshToken       string     `json:"-"`
	TokenExpiresAt     time.Time  `json:"-"`
	RevokedAt          *time.Time `json:"revoked_at,omitempty"`
	InsertedAt         time.Time  `json:"inserted_at"`
}

var (
	GrantMissingScopeError  = errors.New("missing requested scope")
	InvalidAccessTokenError = errors.New("invalid access token")
)

type AuthorizerAPICall struct {
	Method    string
	RoutePath string
	Query     *string
	Body      *string
}

/*
Authorizer authorizes an OAuth2 Access Token and returns either a *Grant if valid,
or an error if invalid for one reason or another-- use errors.Is(...) to compare to
the errors in this package.

You are welcome to leave call as nil-- fill it only if your authorization includes an API
call-- it will be inserted into the oauth2_api_calls table.
*/
func Authorizer(accessToken string, scopes []Scope, call *AuthorizerAPICall) (*Grant, error) {
	grant, grantScopes, err := oauth2.GetGrantAndScopesByAccessToken(util.DB, &oauth2.GetGrantAndScopesByAccessTokenRequest{
		AccessToken:              accessToken,
		AllowInactiveCredentials: false,
		AllowRevokedGrants:       false,
		AllowExpiredGrants:       false,
	})
	if err != nil {
		return nil, fmt.Errorf("error getting grant or listing grant scopes in oauth2 authorizer: %w", err)
	}

	if grant == nil {
		return nil, InvalidAccessTokenError
	}

	scopesAreOK := true

	// go through requested scope
	for _, s := range scopes {
		scopeIsOK := false

		// see if this scope came back from the database
		for _, gs := range *grantScopes {
			if gs == string(s) {
				// if so, mark as such and move to the next scope
				scopeIsOK = true
				break
			}
		}

		// if not, stop
		if !scopeIsOK {
			scopesAreOK = false
			break
		}
	}

	if !scopesAreOK {
		return nil, GrantMissingScopeError
	}

	if call != nil {
		go func(grantID uint64, c *AuthorizerAPICall) {
			err := oauth2.InsertAPICall(util.DB, &oauth2.InsertAPICallRequest{
				GrantID:   grantID,
				RoutePath: c.RoutePath,
				Method:    c.Method,
				Query:     c.Query,
				Body:      c.Body,
			})
			if err != nil {
				// it's not the biggest deal if it errors out
				util.HandleError(fmt.Errorf("error inserting oauth2 api call: %w", err))
				return
			}
		}(grant.ID, call)
	}

	return &Grant{
		ID:                 grant.ID,
		UserID:             grant.UserID,
		OAuth2CredentialID: grant.OAuth2CredentialID,
		RedirectURIID:      grant.RedirectURIID,
		AccessToken:        grant.AccessToken,
		RefreshToken:       grant.RefreshToken,
		TokenExpiresAt:     grant.TokenExpiresAt,
		RevokedAt:          &grant.RevokedAt,
		InsertedAt:         grant.InsertedAt,
	}, nil
}
