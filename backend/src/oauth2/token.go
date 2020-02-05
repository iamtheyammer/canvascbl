package oauth2

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/oauth2"
	"github.com/iamtheyammer/canvascbl/backend/src/middlewares"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type grantType string

const (
	grantTypeAuthorizationCode = grantType("authorization_code")
	grantTypeRefreshToken      = grantType("refresh_token")
)

type tokenHandlerResponse struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	User         struct {
		UserID uint64 `json:"id,omitempty"`
	} `json:"user,omitempty"`
	ExpiresAt string `json:"expires_at,omitempty"`
}

func TokenHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	q := r.URL.Query()

	gt := grantType(q.Get("grant_type"))
	if len(gt) < 1 {
		util.SendBadRequest(w, "missing grant_type as query param")
		return
	} else if gt != grantTypeAuthorizationCode && gt != grantTypeRefreshToken {
		util.SendBadRequest(w, "invalid grant_type as query param")
		return
	}

	clientID := q.Get("client_id")
	if len(clientID) < 1 {
		util.SendBadRequest(w, "missing client_id as query param")
		return
	} else if !util.ValidateUUIDString(clientID) {
		util.SendBadRequest(w, "invalid client_id")
		return
	}

	clientSecret := q.Get("client_secret")
	if len(clientSecret) < 1 {
		util.SendBadRequest(w, "missing client_secret as query param")
		return
	}

	qRedirectURI := q.Get("redirect_uri")
	if len(qRedirectURI) < 1 {
		util.SendBadRequest(w, "missing redirect_uri as query param")
		return
	}

	qCode := q.Get("code")
	qRefreshToken := q.Get("refresh_token")

	// just validation here
	switch gt {
	case grantTypeAuthorizationCode:

		if len(qCode) < 1 {
			util.SendBadRequest(w, "missing code as query param")
			return
		} else if !util.ValidateUUIDString(qCode) {
			util.SendBadRequest(w, "invalid code")
			return
		}
	case grantTypeRefreshToken:
		if len(qRefreshToken) < 1 {
			util.SendBadRequest(w, "missing refresh_token as query param")
			return
		} else if !util.ValidateUUIDString(qRefreshToken) {
			util.SendBadRequest(w, "invalid refresh_token")
			return
		}
	}

	var (
		wg         = sync.WaitGroup{}
		mutex      = sync.Mutex{}
		errs       []error
		credential oauth2.Credential
		code       oauth2.Code
		grant      oauth2.Grant
	)

	// credential
	wg.Add(1)
	go func(id string, secret string) {
		defer wg.Done()

		dbC, dbErr := oauth2.GetCredential(util.DB, &oauth2.ListCredentialsRequest{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			IsActive:     true,
		})
		if dbErr != nil {
			mutex.Lock()
			errs = append(errs, fmt.Errorf("error getting oauth2 credential in token handler: %w", dbErr))
			mutex.Unlock()
			return
		}

		if dbC == nil {
			return
		}

		credential = *dbC
	}(clientID, clientSecret)

	// code
	if gt == grantTypeAuthorizationCode {
		wg.Add(1)
		go func(c string) {
			defer wg.Done()

			dbC, dbErr := oauth2.GetCode(util.DB, &oauth2.ListCodesRequest{Code: qCode})
			if dbErr != nil {
				mutex.Lock()
				errs = append(errs, fmt.Errorf("error getting oauth2 code in token handler: %w", dbErr))
				mutex.Unlock()
				return
			}

			if dbC == nil {
				return
			}

			code = *dbC
		}(qCode)
	}

	if gt == grantTypeRefreshToken {
		wg.Add(1)
		go func(rt string) {
			defer wg.Done()

			dbG, dbErr := oauth2.GetGrant(
				util.DB,
				&oauth2.ListGrantsRequest{
					RefreshToken:       qRefreshToken,
					AllowExpiredTokens: true},
			)
			if dbErr != nil {
				mutex.Lock()
				errs = append(errs, fmt.Errorf("error getting oauth2 grant in token handler: %w", dbErr))
				mutex.Unlock()
				return
			}

			if dbG == nil {
				return
			}

			grant = *dbG
		}(qRefreshToken)
	}

	wg.Wait()

	if len(errs) > 0 {
		for _, e := range errs {
			util.HandleError(e)
		}

		util.SendInternalServerError(w)
		return
	}

	if credential.ID < 1 {
		util.SendUnauthorized(w, "invalid client id/secret")
		return
	}

	if gt == grantTypeAuthorizationCode {
		if code.ID < 1 {
			util.SendUnauthorized(w, "invalid code")
			return
		}

		if code.ExpiresAt.Before(time.Now()) {
			util.SendUnauthorized(w, "expired code, restart oauth2 flow")
			return
		}
	}

	if gt == grantTypeRefreshToken {
		if grant.ID < 1 {
			util.SendUnauthorized(w, "invalid refresh_token")
			return
		}

		if !grant.RevokedAt.IsZero() {
			// this is ON PURPOSE to not leak that a token has been revoked.
			util.SendUnauthorized(w, "invalid refresh_token")
		}
	}

	// now we can get the redirect uri
	rURIIsOK, rURIID, err := oauth2.RedirectURIIsValidForClientID(
		util.DB,
		&oauth2.RedirectURIIsValidForClientIDRequest{
			RedirectURI: qRedirectURI,
			ClientID:    clientID,
		})
	if err != nil {
		util.HandleError(fmt.Errorf("error getting redirect uri: %w", err))
		util.SendInternalServerError(w)
		return
	}

	if !rURIIsOK {
		util.SendBadRequest(w, "invalid redirect_uri")
		return
	}

	// check redirect uri
	if (gt == grantTypeAuthorizationCode && *rURIID != code.RedirectURIID) ||
		(gt == grantTypeRefreshToken && *rURIID != grant.RedirectURIID) {
		util.SendBadRequest(w, "redirect_uri does not match original")
		return
	}

	// we can now take final action: inserting or updating the grant
	if gt == grantTypeAuthorizationCode {
		trx, err := util.DB.Begin()
		if err != nil {
			util.HandleError(fmt.Errorf("error beginning trx at final action: %w", err))
			util.SendInternalServerError(w)
			return
		}

		gReq := oauth2.InsertOAuth2GrantRequest{
			UserID:             code.UserID,
			OAuth2CredentialID: credential.ID,
			RedirectURIID:      *rURIID,
			OAuth2CodeID:       code.ID,
		}

		purpose := q.Get("purpose")
		if len(purpose) > 0 {
			gReq.Purpose = &purpose
		}

		g, err := oauth2.InsertOAuth2Grant(trx, &gReq)
		if err != nil {
			util.HandleError(fmt.Errorf("error inserting oauth2 grant: %w", err))

			rollbackErr := trx.Rollback()
			if rollbackErr != nil {
				util.HandleError(fmt.Errorf("error rolling back trx at final action: %w", rollbackErr))
			}

			util.SendInternalServerError(w)
			return
		}

		err = oauth2.UpdateCode(trx, &oauth2.UpdateCodeRequest{
			Where: oauth2.ListCodesRequest{ID: code.ID},
			Set:   oauth2.InsertOAuth2CodeRequest{Used: true},
		})
		if err != nil {
			util.HandleError(fmt.Errorf("error updating oauth2 code: %w", err))

			rollbackErr := trx.Rollback()
			if rollbackErr != nil {
				util.HandleError(fmt.Errorf("error rolling back trx at update code: %w", rollbackErr))
			}

			util.SendInternalServerError(w)
			return
		}

		ret, err := json.Marshal(&tokenHandlerResponse{
			AccessToken:  g.AccessToken,
			RefreshToken: g.RefreshToken,
			User: struct {
				UserID uint64 `json:"id,omitempty"`
			}{
				UserID: g.UserID,
			},
			ExpiresAt: g.TokenExpiresAt.Format(time.RFC3339),
		})
		if err != nil {
			util.HandleError(fmt.Errorf("error marshaling token handler response: %w", err))

			rollbackErr := trx.Rollback()
			if rollbackErr != nil {
				util.HandleError(fmt.Errorf("error rolling back trx at final action: %w", rollbackErr))
			}

			util.SendInternalServerError(w)
			return
		}

		err = trx.Commit()
		if err != nil {
			util.HandleError(fmt.Errorf("error committing transaction at final action: %w", err))
			util.SendInternalServerError(w)
			return
		}

		util.SendJSONResponse(w, ret)
		return
	}

	if gt == grantTypeRefreshToken {
		cycledGrant, err := oauth2.CycleGrantAccessToken(util.DB, grant.RefreshToken)
		if err != nil {
			util.HandleError(fmt.Errorf("error cycling grant access token: %w", err))
			util.SendInternalServerError(w)
			return
		}

		ret, err := json.Marshal(&tokenHandlerResponse{
			AccessToken: cycledGrant.AccessToken,
			User: struct {
				UserID uint64 `json:"id,omitempty"`
			}{
				grant.UserID,
			},
			ExpiresAt: cycledGrant.TokenExpiresAt.Format(time.RFC3339),
		})
		if err != nil {
			util.HandleError(fmt.Errorf("error marshaling tokenHandlerResponse for refresh token: %w", err))
			util.SendInternalServerError(w)
			return
		}

		util.SendJSONResponse(w, ret)
		return
	}

	// something went wrong.
	util.SendInternalServerError(w)
	return
}

func DeleteTokenHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var revokeReq *oauth2.RevokeGrantRequest

	sess := middlewares.Session(w, r, false)
	if sess == nil {
		at, tokenIsOK := middlewares.Bearer(w, r, true)
		if !tokenIsOK {
			util.SendUnauthorized(w, "invalid access token")
			return
		}
		g, err := Authorizer(at, []Scope{}, nil)
		if err != nil {
			if errors.Is(err, InvalidAccessTokenError) {
				util.SendUnauthorized(w, "invalid access token")
				return
			}

			util.HandleError(fmt.Errorf("error in delete token handler Authorizer: %w", err))
			util.SendInternalServerError(w)
			return
		}

		revokeReq = &oauth2.RevokeGrantRequest{
			ID:                 g.ID,
			UserID:             g.UserID,
			OAuth2CredentialID: g.OAuth2CredentialID,
			RedirectURIID:      g.RedirectURIID,
			AccessToken:        g.AccessToken,
			RefreshToken:       g.RefreshToken,
		}
	} else {
		revID := r.URL.Query().Get("token_id")
		if len(revID) < 1 || !util.ValidateIntegerString(revID) {
			util.SendBadRequest(w, "missing token_id as query param")
			return
		}

		iRevID, err := strconv.Atoi(revID)
		if err != nil {
			util.HandleError(fmt.Errorf("error converting to revoke id to an int: %w", err))
			util.SendInternalServerError(w)
			return
		}

		revokeReq = &oauth2.RevokeGrantRequest{
			ID:     uint64(iRevID),
			UserID: sess.UserID,
		}
	}

	err := oauth2.RevokeGrant(util.DB, revokeReq)
	if err != nil {
		util.HandleError(fmt.Errorf("error revoking oauth2 grant: %w", err))
		util.SendInternalServerError(w)
		return
	}

	util.SendNoContent(w)
	return
}
