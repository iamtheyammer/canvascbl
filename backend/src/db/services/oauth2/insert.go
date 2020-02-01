package oauth2

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
)

type InsertOAuth2CodeRequest struct {
	UserID             uint64
	OAuth2CredentialID uint64
	RedirectURIID      uint64
	Used               bool
	ScopeIDs           []uint64
}

type InsertOAuth2GrantRequest struct {
	UserID             uint64
	OAuth2CredentialID uint64
	RedirectURIID      uint64
	OAuth2CodeID       uint64
}

type InsertAPICallRequest struct {
	GrantID   uint64
	RouteID   uint64
	RoutePath string
	Query     *string
	Body      *string
}

/*
InsertOAuth2Code inserts an OAuth2 Code into oauth2_codes, along with its scopes
in oauth2_scope_grants. Should be used with a transaction.

Note that InsertOAuth2Code DOES NOT verify that the requested OAuth2 Credential
has all the scopes requested-- you must do this yourself with
GetOAuth2CredentialScopes.
*/
func InsertOAuth2Code(db services.DB, req *InsertOAuth2CodeRequest) (*Code, error) {
	query, args, err := util.Sq.
		Insert("oauth2_codes").
		SetMap(map[string]interface{}{
			"user_id":              req.UserID,
			"oauth2_credential_id": req.OAuth2CredentialID,
			"redirect_uri_id":      req.RedirectURIID,
			"used":                 req.Used,
		}).
		Suffix("RETURNING id, code, consent_code, expires_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building insert oauth2 code sql: %w", err)
	}

	row := db.QueryRow(query, args...)

	var c Code
	err = row.Scan(&c.ID, &c.Code, &c.ConsentCode, &c.ExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("error scanning insert oauth2 code sql: %w", err)
	}

	q := util.Sq.
		Insert("oauth2_scope_grants").
		Columns(
			"oauth2_credential_id",
			"scope_id",
			"oauth2_code_id",
		)

	for _, sID := range req.ScopeIDs {
		q = q.Values(req.OAuth2CredentialID, sID, c.ID)
	}

	query, args, err = q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building insert oauth2 code insert scopes sql: %w", err)
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing insert oauth2 code insert scopes sql: %w", err)
	}

	return &c, nil
}

/*
InsertOAuth2Grant inserts a grant along with updating the scope grants with the grant ID.
Should be used in a transaction.
*/
func InsertOAuth2Grant(db services.DB, req *InsertOAuth2GrantRequest) (*Grant, error) {
	query, args, err := util.Sq.
		Insert("oauth2_grants").
		SetMap(map[string]interface{}{
			"user_id":              req.UserID,
			"oauth2_credential_id": req.OAuth2CredentialID,
			"redirect_uri_id":      req.RedirectURIID,
		}).
		Suffix("RETURNING id, access_token, refresh_token, token_expires_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building insert oauth2 grant sql: %w", err)
	}

	row := db.QueryRow(query, args...)

	g := Grant{
		UserID:             req.UserID,
		OAuth2CredentialID: req.OAuth2CredentialID,
	}
	err = row.Scan(&g.ID, &g.AccessToken, &g.RefreshToken, &g.TokenExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("error executing insert oauth2 grant sql: %w", err)
	}

	query, args, err = util.Sq.
		Update("oauth2_scope_grants").
		Set("oauth2_grant_id", g.ID).
		Where(sq.Eq{"oauth2_code_id": req.OAuth2CodeID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building insert oauth2 grant update sql: %w", err)
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing insert oauth2 grant update sql: %w", err)
	}

	return &g, nil
}

func InsertAPICall(db services.DB, req *InsertAPICallRequest) error {
	q := util.Sq.
		Insert("oauth2_api_calls").
		Columns(
			"oauth2_grant_id",
			"api_route_id",
			"query",
			"body",
		)

	var routeID interface{}
	if req.RouteID > 0 {
		routeID = req.RouteID
	} else if len(req.RoutePath) > 0 {
		//q = q.Prefix("WITH rid AS (SELECT id FROM api_routes WHERE path = ?)", req.RoutePath)
		routeID = sq.Expr("(SELECT id FROM api_routes WHERE path = ?)", req.RoutePath)
	}

	q = q.Values(req.GrantID, routeID, req.Query, req.Body)

	query, args, err := q.ToSql()
	if err != nil {
		return fmt.Errorf("error building insert api call sql: %w", err)
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("error executing insert api call sql: %w", err)
	}

	return nil
}
