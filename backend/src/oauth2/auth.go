package oauth2

import (
	"fmt"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/oauth2"
	"github.com/iamtheyammer/canvascbl/backend/src/env"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

func AuthHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	q := r.URL.Query()

	clientID := q.Get("client_id")
	if len(clientID) < 1 {
		util.SendBadRequest(w, "missing client_id as query param")
		return
	} else if !util.ValidateUUIDString(clientID) {
		util.SendBadRequest(w, "invalid client_id as query param")
		return
	}

	// only supported value is "code"
	if q.Get("response_type") != "code" {
		util.SendBadRequest(w, "missing/invalid response_type as query param")
		return
	}

	redirectURI := q.Get("redirect_uri")
	if len(redirectURI) < 1 {
		util.SendBadRequest(w, "missing redirect_uri as query param")
		return
	}

	qScopes := q.Get("scope")
	scopes := strings.Split(qScopes, " ")
	if qScopes == "" {
		util.SendBadRequest(w, "missing scope as query param")
		return
	} else if ok, invScope := ValidateScopes(scopes); !ok {
		util.SendBadRequest(w, "unknown scope: "+*invScope)
		return
	}

	// ensure all requested scopes are ok and test redirect uri
	// doing this with goroutines to speed it up
	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}
	var err error

	var (
		redirectURIIsOK, scopesAreOK bool
		credentialID, redirectURIID  uint64
		scopeIDs                     []uint64
	)

	// redirect URI; using params makes them concurrent-safe without mutex use
	wg.Add(1)
	go func(cID string, rURI string) {
		defer wg.Done()

		ok, rURIID, rURIErr := oauth2.RedirectURIIsValidForClientID(util.DB, &oauth2.RedirectURIIsValidForClientIDRequest{
			ClientID:    cID,
			RedirectURI: rURI,
		})
		if rURIErr != nil {
			mutex.Lock()
			err = fmt.Errorf("error figuring out whether redirect uri %s is valid for client id %s: %w", cID, rURI, rURIErr)
			mutex.Unlock()
			return
		}

		if rURIID != nil {
			redirectURIIsOK = ok
			redirectURIID = *rURIID
		}
	}(clientID, redirectURI)

	// scopes
	wg.Add(1)
	go func(cID string) {
		defer wg.Done()

		dbS, credID, ssErr := oauth2.GetOAuth2CredentialScopes(util.DB, &oauth2.GetOAuth2CredentialScopesRequest{
			CredentialClientID: cID,
		})
		if ssErr != nil {
			mutex.Lock()
			err = fmt.Errorf("error getting oauth2 credential scopes for client id %s: %w", cID, ssErr)
			mutex.Unlock()
			return
		}

		credentialID = *credID

		// check scopes
		for _, s := range scopes {
			found := false
			for _, dbS := range *dbS {
				if s == dbS.ShortName {
					found = true
					scopeIDs = append(scopeIDs, dbS.ID)
					// why continue if we found it?
					break
				}
			}

			if !found {
				scopesAreOK = false
				return
			}
		}

		scopesAreOK = true
		return
	}(clientID)

	wg.Wait()

	if err != nil {
		util.SendInternalServerError(w)
		return
	}

	if credentialID < 1 {
		util.SendBadRequest(w, "invalid client_id")
		return
	}

	if !redirectURIIsOK {
		util.SendBadRequest(w, "invalid redirect uri")
		return
	}

	if !scopesAreOK {
		util.SendBadRequest(w, "unauthorized scope")
		return
	}

	trx, err := util.DB.Begin()
	if err != nil {
		util.HandleError(fmt.Errorf("error beginning db trx for oauth2 requesthandler: %w", err))
		util.SendInternalServerError(w)
		return
	}

	c, err := oauth2.InsertOAuth2Code(trx, &oauth2.InsertOAuth2CodeRequest{
		OAuth2CredentialID: credentialID,
		RedirectURIID:      redirectURIID,
		ScopeIDs:           scopeIDs,
	})
	if err != nil {
		rollbackErr := trx.Rollback()
		if rollbackErr != nil {
			util.HandleError(fmt.Errorf("error rolling back oauth2 request transaction at insert oauth2 code: %w", rollbackErr))
		}

		util.HandleError(fmt.Errorf("error inserting oauth2 code: %w", err))
		util.SendInternalServerError(w)
		return
	}

	err = trx.Commit()
	if err != nil {
		util.HandleError(fmt.Errorf("error committing oauth2 request transaction: %w", err))
		rollbackErr := trx.Rollback()
		if rollbackErr != nil {
			util.HandleError(fmt.Errorf("error rolling back oauth2 request transaction at commit: %w", err))
		}

		util.SendInternalServerError(w)
		return
	}

	rq := url.Values{}
	rq.Add("type", "oauth2")
	// we know that all scopes are OK, so we can return those short names
	rq.Add("scope", qScopes)
	rq.Add("consent_code", c.ConsentCode)

	util.SendRedirect(w, env.OAuth2ConsentURL+"?"+rq.Encode())
	return
}
