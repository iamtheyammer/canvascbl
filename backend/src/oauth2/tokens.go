package oauth2

import (
	"encoding/json"
	"fmt"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/oauth2"
	"github.com/iamtheyammer/canvascbl/backend/src/middlewares"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type Credential struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

func TokensHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sess := middlewares.Session(w, r, true)
	if sess == nil {
		return
	}

	grants, err := oauth2.ListGrants(util.DB, &oauth2.ListGrantsRequest{
		UserID:                   sess.UserID,
		AllowExpiredTokens:       true,
		AllowInactiveCredentials: true,
		AllowRevoked:             false,
	})
	if err != nil {
		util.HandleError(fmt.Errorf("error listing oauth2 grants: %w", err))
		util.SendInternalServerError(w)
		return
	}

	var (
		gs    []Grant
		mCIDs = make(map[uint64]struct{})
	)

	for _, g := range *grants {
		gs = append(gs, Grant{
			ID:                 g.ID,
			UserID:             g.UserID,
			OAuth2CredentialID: g.OAuth2CredentialID,
			RedirectURIID:      g.RedirectURIID,
			AccessToken:        g.AccessToken,
			RefreshToken:       g.RefreshToken,
			TokenExpiresAt:     g.TokenExpiresAt,
			RevokedAt:          g.RevokedAt,
			InsertedAt:         g.InsertedAt,
		})
		mCIDs[g.OAuth2CredentialID] = struct{}{}
	}

	// credentials
	var cIDs []uint64
	for cID := range mCIDs {
		cIDs = append(cIDs, cID)
	}

	creds, err := oauth2.ListCredentials(util.DB, &oauth2.ListCredentialsRequest{IDs: cIDs})
	if err != nil {
		util.HandleError(fmt.Errorf("error listing oauth2 credentials"))
		util.SendInternalServerError(w)
		return
	}

	var cs []Credential
	for _, c := range *creds {
		cs = append(cs, Credential{
			ID:   c.ID,
			Name: c.Name,
		})
	}

	jGs, err := json.Marshal(&struct {
		Credentials []Credential `json:"credentials"`
		Grants      []Grant      `json:"grants"`
	}{
		Credentials: cs,
		Grants:      gs,
	})
	if err != nil {
		util.HandleError(fmt.Errorf("error marshaling json tokens response: %w", err))
		util.SendInternalServerError(w)
		return
	}

	util.SendJSONResponse(w, jGs)
	return
}
