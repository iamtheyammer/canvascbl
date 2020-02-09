package oauth2

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
)

type UpdateCodeRequest struct {
	Where ListCodesRequest
	Set   InsertOAuth2CodeRequest
}

//type UpdateGrantRequestSet struct {
//	UserID             uint64
//	OAuth2CredentialID uint64
//	RedirectURIID      uint64
//
//	CycleAccessToken bool
//}

//type UpdateGrantRequest struct {
//	Where Grant
//	Set   UpdateGrantRequestSet
//}

// CycleGrantAccessToken resets access_token and token_expires_at
func CycleGrantAccessToken(db services.DB, refreshToken string) (*Grant, error) {
	// yes, I know I shouldn't be writing SQL but there's no other way to do it that
	// allows setting to DEFAULT
	args := []interface{}{refreshToken}
	query, err := util.PlaceholderFormat.ReplacePlaceholders(
		"UPDATE oauth2_grants SET token_expires_at = DEFAULT, access_token = DEFAULT " +
			"WHERE refresh_token = ? " +
			"RETURNING id, user_id, oauth2_credential_id, redirect_uri_id, access_token, refresh_token, " +
			"token_expires_at, inserted_at",
	)
	if err != nil {
		return nil, fmt.Errorf("error replacing placeholders in cycle grant access token sql: %w", err)
	}

	row := db.QueryRow(query, args...)

	var g Grant
	err = row.Scan(
		&g.ID,
		&g.UserID,
		&g.OAuth2CredentialID,
		&g.RedirectURIID,
		&g.AccessToken,
		&g.RefreshToken,
		&g.TokenExpiresAt,
		&g.InsertedAt,
	)
	if err != nil {
		// there should always be a row-- otherwise there's other errors
		return nil, fmt.Errorf("error executing cycle grant access token sql: %w", err)
	}

	return &g, nil
}

//func UpdateGrants(db services.DB, req *UpdateGrantRequest) (*Grant, error) {
//	q := util.Sq.
//		Update("oauth2_grants")
//
//	if req.Where.ID > 0 {
//		q = q.Where(sq.Eq{"id": req.Where.ID})
//	}
//
//	if req.Where.UserID > 0 {
//		q = q.Where(sq.Eq{"user_id": req.Where.UserID})
//	}
//
//	if req.Where.OAuth2CredentialID > 0 {
//		q = q.Where(sq.Eq{"oauth2_credential_id": req.Where.OAuth2CredentialID})
//	}
//
//	if len(req.Where.AccessToken) > 0 {
//		q = q.Where(sq.Eq{"access_token": req.Where.AccessToken})
//	}
//
//	if len(req.Where.RefreshToken) > 0 {
//		q = q.Where(sq.Eq{"refresh_token": req.Where.RefreshToken})
//	}
//
//	if req.Set.UserID > 0 {
//		q = q.Set("user_id", req.Set.UserID)
//	}
//
//	if req.Set.OAuth2CredentialID > 0 {
//		q = q.Set("oauth2_credential_id", req.Set.OAuth2CredentialID)
//	}
//
//	if req.Set.RedirectURIID > 0 {
//		q = q.Set("redirect_uri_id", req.Set.RedirectURIID)
//	}
//
//	query, args, err := q.ToSql()
//	if err != nil {
//		return nil, fmt.Errorf("error updating grants: %w", err)
//	}
//}

func UpdateCode(db services.DB, req *UpdateCodeRequest) error {
	q := util.Sq.
		Update("oauth2_codes").
		Where(sq.Eq{"used": req.Where.AllowUsed})

	if req.Where.ID > 0 {
		q = q.Where(sq.Eq{"id": req.Where.ID})
	}

	if req.Where.UserID > 0 {
		q = q.Where(sq.Eq{"user_id": req.Where.UserID})
	}

	if req.Where.OAuth2CredentialID > 0 {
		q = q.Where(sq.Eq{"oauth2_credential_id": req.Where.OAuth2CredentialID})
	}

	if req.Where.RedirectURIID > 0 {
		q = q.Where(sq.Eq{"redirect_uri_id": req.Where.RedirectURIID})
	}

	if len(req.Where.Code) > 0 {
		q = q.Where(sq.Eq{"code": req.Where.Code})
	}

	if len(req.Where.ConsentCode) > 0 {
		q = q.Where(sq.Eq{"consent_code": req.Where.ConsentCode})
	}

	if req.Set.UserID != nil {
		q = q.Set("user_id", req.Set.UserID)
	}

	if req.Set.OAuth2CredentialID > 0 {
		q = q.Set("oauth2_credential_id", req.Set.OAuth2CredentialID)
	}

	if req.Set.RedirectURIID > 0 {
		q = q.Set("redirect_uri_id", req.Set.RedirectURIID)
	}

	if req.Set.Used {
		q = q.Set("used", req.Set.Used)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return fmt.Errorf("error building update oauth2 code sql: %w", err)
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("error executing update oauth2 code sql: %w", err)
	}

	return nil
}
