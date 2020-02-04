package oauth2

import (
	"database/sql"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"time"
)

type Code struct {
	ID                 uint64
	UserID             uint64
	OAuth2CredentialID uint64
	RedirectURIID      uint64
	Code               string
	ConsentCode        string
	Used               bool
	ExpiresAt          time.Time
	InsertedAt         time.Time
}

type Grant struct {
	ID                 uint64
	UserID             uint64
	OAuth2CredentialID uint64
	RedirectURIID      uint64
	AccessToken        string
	RefreshToken       string
	Purpose            string
	TokenExpiresAt     time.Time
	RevokedAt          time.Time
	InsertedAt         time.Time
}

type Scope struct {
	ID          uint64
	ShortName   string
	Description string
	//InsertedAt string
}

type RedirectURI struct {
	ID                 uint64
	OAuth2CredentialID uint64
	RedirectURI        string
	InsertedAt         time.Time
}

type Credential struct {
	ID           uint64
	Name         string
	OwnerUserID  uint64
	ClientID     string
	ClientSecret string
	IsActive     bool
	InsertedAt   time.Time
}

type GetOAuth2CredentialScopesRequest struct {
	// Client ID of the requested credential
	CredentialClientID string

	// Whether you're ok with inactive credentials
	AllowInactiveCredentials bool
}

type RedirectURIIsValidForClientIDRequest struct {
	ClientID                 string
	RedirectURI              string
	AllowInactiveCredentials bool
}

type ListCodesRequest struct {
	ID                 uint64
	UserID             uint64
	OAuth2CredentialID uint64
	RedirectURIID      uint64
	Code               string
	ConsentCode        string

	// used, for update. true would set used to true, false for false
	AllowUsed bool
}

type ListRedirectURIsRequest struct {
	ID                 uint64
	OAuth2CredentialID uint64
	RedirectURI        string
}

type ListCredentialsRequest struct {
	ID           uint64
	IDs          []uint64
	OwnerUserID  uint64
	ClientID     string
	ClientSecret string
	IsActive     bool
}

type ListGrantsRequest struct {
	ID                 uint64
	UserID             uint64
	OAuth2CredentialID uint64
	AccessToken        string
	RefreshToken       string

	AllowExpiredTokens       bool
	AllowRevoked             bool
	AllowInactiveCredentials bool
}

type ListGrantScopesRequest struct {
	ID                 uint64
	UserID             uint64
	OAuth2CredentialID uint64
	RedirectURIID      uint64
	AccessToken        string
	RefreshToken       string

	AllowInactiveCredentials bool
	AllowRevoked             bool
	Limit                    uint64
}

/*
GetOAuth2Credential scopes returns a list of Scopes on the Credential (if one exists) along
with the credential's ID (again, if one exists).
*/
func GetOAuth2CredentialScopes(db services.DB, req *GetOAuth2CredentialScopesRequest) (*[]Scope, *uint64, error) {
	q := util.Sq.
		Select("oauth2_scopes.id", "oauth2_scopes.short_name", "oauth2_credentials.id").
		From("oauth2_credentials_scopes").
		Join("oauth2_credentials ON oauth2_credentials_scopes.oauth2_credential_id = oauth2_credentials.id").
		Join("oauth2_scopes ON oauth2_credentials_scopes.scope_id = oauth2_scopes.id").
		Where(sq.Eq{"oauth2_credentials.client_id": req.CredentialClientID})

	if !req.AllowInactiveCredentials {
		q = q.Where(sq.Eq{"oauth2_credentials.is_active": true})
	}

	query, args, err := q.ToSql()
	if err != nil {
		return nil, nil, fmt.Errorf("error building get oauth2 credential scopes query: %w", err)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("error executing get oauth2 credential scopes query: %w", err)
	}

	defer rows.Close()

	var (
		scopes []Scope
		cID    uint64
	)
	for rows.Next() {
		var s Scope
		err = rows.Scan(&s.ID, &s.ShortName, &cID)
		if err != nil {
			return nil, nil, fmt.Errorf("error scanning get oauth2 credential scopes query: %w", err)
		}

		scopes = append(scopes, s)
	}

	return &scopes, &cID, nil
}

/*
RedirectURIIsValidForClientID returns a bool (and an error) indicating if the specified Redirect URI
is a valid one for the Client ID submitted.

You can use it like:

ok, err := oauth2.RedirectURIIsValidForClientID(db, req)
*/
func RedirectURIIsValidForClientID(db services.DB, req *RedirectURIIsValidForClientIDRequest) (bool, *uint64, error) {
	query, args, err := util.Sq.
		Select("oauth2_credentials_redirect_uris.id").
		From("oauth2_credentials_redirect_uris").
		Join("oauth2_credentials ON oauth2_credentials_redirect_uris.oauth2_credential_id = oauth2_credentials.id").
		Where(sq.Eq{"oauth2_credentials_redirect_uris.redirect_uri": req.RedirectURI}).
		Where(sq.Eq{"oauth2_credentials.client_id": req.ClientID}).
		Where(sq.Eq{"oauth2_credentials.is_active": !req.AllowInactiveCredentials}).
		ToSql()
	if err != nil {
		return false, nil, fmt.Errorf("error building redirect uri is valid for client id sql: %w", err)
	}

	row := db.QueryRow(query, args...)

	var redirectURIID uint64
	err = row.Scan(&redirectURIID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil, nil
		}

		return false, nil, fmt.Errorf("error scanning redirect uri is valid for client id sql: %w", err)
	}

	return true, &redirectURIID, nil
}

func GetCode(db services.DB, req *ListCodesRequest) (*Code, error) {
	q := util.Sq.
		Select(
			"id",
			"user_id",
			"oauth2_credential_id",
			"redirect_uri_id",
			"code",
			"consent_code",
			"used",
			"expires_at",
			"inserted_at",
		).
		From("oauth2_codes")

	if req.ID > 0 {
		q = q.Where(sq.Eq{"id": req.ID})
	}

	if req.UserID > 0 {
		q = q.Where(sq.Eq{"user_id": req.UserID})
	}

	if req.OAuth2CredentialID > 0 {
		q = q.Where(sq.Eq{"oauth2_credential_id": req.OAuth2CredentialID})
	}

	if req.RedirectURIID > 0 {
		q = q.Where(sq.Eq{"redirect_uri_id": req.RedirectURIID})
	}

	if len(req.Code) > 0 {
		q = q.Where(sq.Eq{"code": req.Code})
	}

	if len(req.ConsentCode) > 0 {
		q = q.Where(sq.Eq{"consent_code": req.ConsentCode})
	}

	if !req.AllowUsed {
		q = q.Where(sq.Eq{"used": req.AllowUsed})
	}

	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building get oauth2 code sql: %w", err)
	}

	row := db.QueryRow(query, args...)

	var c Code
	err = row.Scan(
		&c.ID,
		&c.UserID,
		&c.OAuth2CredentialID,
		&c.RedirectURIID,
		&c.Code,
		&c.ConsentCode,
		&c.Used,
		&c.ExpiresAt,
		&c.InsertedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error scanning get oauth2 code sql: %w", err)
	}

	return &c, nil
}

func GetRedirectURI(db services.DB, req *ListRedirectURIsRequest) (*RedirectURI, error) {
	q := util.Sq.
		Select(
			"id",
			"oauth2_credential_id",
			"redirect_uri",
			"inserted_at",
		).
		From("oauth2_credentials_redirect_uris")

	if req.ID > 0 {
		q = q.Where(sq.Eq{"id": req.ID})
	}

	if req.OAuth2CredentialID > 0 {
		q = q.Where(sq.Eq{"oauth2_credential_id": req.OAuth2CredentialID})
	}

	if len(req.RedirectURI) > 0 {
		q = q.Where(sq.Eq{"redirect_uri": req.RedirectURI})
	}

	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building get redirect uri sql: %w", err)
	}

	row := db.QueryRow(query, args...)

	var r RedirectURI
	err = row.Scan(
		&r.ID,
		&r.OAuth2CredentialID,
		&r.RedirectURI,
		&r.InsertedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error scanning get redirect uri row: %w", err)
	}

	return &r, nil
}

func ListCredentials(db services.DB, req *ListCredentialsRequest) (*[]Credential, error) {
	q := util.Sq.
		Select(
			"id",
			"name",
			"owner_user_id",
			"client_id",
			"client_secret",
			"is_active",
			"inserted_at",
		).
		From("oauth2_credentials")

	if req.ID > 0 {
		q = q.Where(sq.Eq{"id": req.ID})
	}

	if len(req.IDs) > 0 {
		q = q.Where(sq.Eq{"id": req.IDs})
	}

	if req.OwnerUserID > 0 {
		q = q.Where(sq.Eq{"owner_user_id": req.OwnerUserID})
	}

	if len(req.ClientID) > 0 {
		q = q.Where(sq.Eq{"client_id": req.ClientID})
	}

	if len(req.ClientSecret) > 0 {
		q = q.Where(sq.Eq{"client_secret": req.ClientSecret})
	}

	if req.IsActive {
		q = q.Where(sq.Eq{"is_active": req.IsActive})
	}

	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building list credentials sql: %w", err)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing list credentials sql: %w", err)
	}

	var cs []Credential
	for rows.Next() {
		var c Credential
		err = rows.Scan(
			&c.ID,
			&c.Name,
			&c.OwnerUserID,
			&c.ClientID,
			&c.ClientSecret,
			&c.IsActive,
			&c.InsertedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning credentials from list credential sql: %w", err)
		}

		cs = append(cs, c)
	}

	return &cs, nil
}

func GetCredential(db services.DB, req *ListCredentialsRequest) (*Credential, error) {
	cs, err := ListCredentials(db, req)
	if err != nil {
		return nil, fmt.Errorf("error listing credntials in get credential: %w", err)
	}

	if len(*cs) < 1 {
		return nil, nil
	}

	return &(*cs)[0], nil
}

func ListGrants(db services.DB, req *ListGrantsRequest) (*[]Grant, error) {
	q := util.Sq.
		Select(
			"oauth2_grants.id",
			"oauth2_grants.user_id",
			"oauth2_grants.purpose",
			"oauth2_grants.oauth2_credential_id",
			"oauth2_grants.redirect_uri_id",
			"oauth2_grants.access_token",
			"oauth2_grants.refresh_token",
			"oauth2_grants.token_expires_at",
			"oauth2_grants.revoked_at",
			"oauth2_grants.inserted_at",
		).
		From("oauth2_grants")

	if req.ID > 0 {
		q = q.Where(sq.Eq{"id": req.ID})
	}

	if req.UserID > 0 {
		q = q.Where(sq.Eq{"user_id": req.UserID})
	}

	if req.OAuth2CredentialID > 0 {
		q = q.Where(sq.Eq{"oauth2_credential_id": req.OAuth2CredentialID})
	}

	if len(req.AccessToken) > 0 {
		q = q.Where(sq.Eq{"access_token": req.AccessToken})
	}

	if len(req.RefreshToken) > 0 {
		q = q.Where(sq.Eq{"refresh_token": req.RefreshToken})
	}

	if !req.AllowRevoked {
		q = q.Where(sq.Eq{"revoked_at": nil})
	}

	if !req.AllowInactiveCredentials {
		q = q.Join("oauth2_credentials ON oauth2_grants.oauth2_credential_id = oauth2_credentials.id").
			Where(sq.Eq{"oauth2_credentials.is_active": true})
	}

	// if we do not allow expired tokens
	if !req.AllowExpiredTokens {
		// then the expiry time needs to be after now
		q = q.Where("token_expires_at > NOW()")
	}

	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building list grants sql: %w", err)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing list grants sql: %w", err)
	}

	var gs []Grant
	for rows.Next() {
		var (
			g         Grant
			purpose   sql.NullString
			revokedAt sql.NullTime
		)
		err = rows.Scan(
			&g.ID,
			&g.UserID,
			&purpose,
			&g.OAuth2CredentialID,
			&g.RedirectURIID,
			&g.AccessToken,
			&g.RefreshToken,
			&g.TokenExpiresAt,
			&revokedAt,
			&g.InsertedAt,
		)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, nil
			}
			return nil, fmt.Errorf("error scanning list grants sql: %w", err)
		}

		if purpose.Valid {
			g.Purpose = purpose.String
		}

		if revokedAt.Valid {
			g.RevokedAt = revokedAt.Time
		}

		gs = append(gs, g)
	}

	return &gs, nil
}

func GetGrant(db services.DB, req *ListGrantsRequest) (*Grant, error) {
	gs, err := ListGrants(db, req)
	if err != nil {
		return nil, fmt.Errorf("error listing grants in get grant: %w", err)
	}

	if len(*gs) < 1 {
		return nil, nil
	}

	return &(*gs)[0], nil
}

func ListGrantScopes(db services.DB, req *ListGrantScopesRequest) (*[]Scope, error) {
	q := util.Sq.
		Select(
			"oauth2_scopes.id",
			"oauth2_scopes.short_name",
			"oauth2_scopes.description",
		).
		From("oauth2_scopes").
		Join("oauth2_scope_grants ON oauth2_scopes.id = oauth2_scope_grants.scope_id").
		LeftJoin("oauth2_grants ON oauth2_scope_grants.oauth2_grant_id = oauth2_grants.id")

	if req.ID > 0 {
		q = q.Where(sq.Eq{"id": req.ID})
	}

	if req.UserID > 0 {
		q = q.Where(sq.Eq{"oauth2_grants.user_id": req.UserID})
	}

	if req.OAuth2CredentialID > 0 {
		q = q.Where(sq.Eq{"oauth2_grants.oauth2_credential_id": req.OAuth2CredentialID})
	}

	if req.RedirectURIID > 0 {
		q = q.Where(sq.Eq{"oauth2_grants.redirect_uri_id": req.RedirectURIID})
	}

	if len(req.AccessToken) > 0 {
		q = q.Where(sq.Eq{"oauth2_grants.access_token": req.AccessToken})
	}

	if len(req.RefreshToken) > 0 {
		q = q.Where(sq.Eq{"oauth2_grants.refresh_token": req.RefreshToken})
	}

	if !req.AllowInactiveCredentials {
		q = q.LeftJoin("oauth2_credentials ON oauth2_scope_grants.oauth2_credential_id = oauth2_credentials.id").
			Where(sq.Eq{"oauth2_credentials.is_active": !req.AllowInactiveCredentials})
	}

	if !req.AllowRevoked {
		q = q.Where(sq.Eq{"oauth2_grants.revoked_at": nil})
	}

	if req.Limit > 0 {
		q = q.Limit(req.Limit)
	} else {
		q = q.Limit(services.DefaultSelectLimit)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building list grant scopes sql: %w", err)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing list grant scopes sql: %w", err)
	}

	defer rows.Close()

	var scopes []Scope
	for rows.Next() {
		var s Scope
		err = rows.Scan(
			&s.ID,
			&s.ShortName,
			&s.Description,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning list grant scopes sql: %w", err)
		}

		scopes = append(scopes, s)
	}

	return &scopes, nil
}
