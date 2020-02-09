package oauth2

import (
	"errors"
	"fmt"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/oauth2"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"sync"
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
	var (
		wg          = sync.WaitGroup{}
		mutex       = sync.Mutex{}
		err         error
		g           oauth2.Grant
		scopesAreOK bool
	)

	// go get grant
	wg.Add(1)
	go func(at string) {
		defer wg.Done()

		gr, grErr := oauth2.GetGrant(util.DB, &oauth2.ListGrantsRequest{
			AccessToken:              at,
			AllowRevoked:             false,
			AllowInactiveCredentials: false,
			AllowExpiredTokens:       false,
		})
		if grErr != nil {
			mutex.Lock()
			err = fmt.Errorf("error listing grants in oauth2 authorizer: %w", grErr)
			mutex.Unlock()
		}

		// careful not to dereference nil
		if gr != nil {
			g = *gr
		}
	}(accessToken)

	// go get scopes
	wg.Add(1)
	go func(at string, reqScopes []Scope) {
		defer wg.Done()

		s, sErr := oauth2.ListGrantScopes(util.DB, &oauth2.ListGrantScopesRequest{
			AccessToken:              at,
			AllowInactiveCredentials: false,
			AllowRevoked:             false,
		})
		if sErr != nil {
			mutex.Lock()
			err = fmt.Errorf("error listing grant scopes in oauth2 authorizer: %w", sErr)
			mutex.Unlock()
		}

		if s == nil {
			// we can return as the lack of a grant will make up for it
			return
		}

		if len(*s) < 1 {
			return
		}

		// go through requested scope
		for _, rs := range reqScopes {
			scopeIsOK := false

			// see if this scope came back from the database
			for _, sc := range *s {
				if sc.ShortName == string(rs) {
					// if so, mark as such and move to the next scope
					scopeIsOK = true
					break
				}
			}

			// if not, stop
			if !scopeIsOK {
				scopesAreOK = false
				return
			}
		}

		scopesAreOK = true
		return
	}(accessToken, scopes)

	wg.Wait()

	if err != nil {
		return nil, fmt.Errorf("error getting grant or listing grant scopes in oauth2 authorizer: %w", err)
	}

	if g.ID < 1 {
		return nil, InvalidAccessTokenError
	}

	if !scopesAreOK {
		return nil, GrantMissingScopeError
	}

	if call != nil {
		go func(grant oauth2.Grant, c *AuthorizerAPICall) {
			err := oauth2.InsertAPICall(util.DB, &oauth2.InsertAPICallRequest{
				GrantID:   grant.ID,
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
		}(g, call)
	}

	return &Grant{
		ID:                 g.ID,
		UserID:             g.UserID,
		OAuth2CredentialID: g.OAuth2CredentialID,
		RedirectURIID:      g.RedirectURIID,
		AccessToken:        g.AccessToken,
		RefreshToken:       g.RefreshToken,
		TokenExpiresAt:     g.TokenExpiresAt,
		RevokedAt:          &g.RevokedAt,
		InsertedAt:         g.InsertedAt,
	}, nil
}
