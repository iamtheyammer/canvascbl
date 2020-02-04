package oauth2

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
)

type RevokeGrantRequest struct {
	ID                 uint64
	UserID             uint64
	OAuth2CredentialID uint64
	RedirectURIID      uint64
	AccessToken        string
	RefreshToken       string
}

func RevokeGrant(db services.DB, req *RevokeGrantRequest) error {
	q := util.Sq.
		Update("oauth2_grants").
		Set("revoked_at", sq.Expr("NOW()"))

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

	if len(req.AccessToken) > 0 {
		q = q.Where(sq.Eq{"access_token": req.AccessToken})
	}

	if len(req.RefreshToken) > 0 {
		q = q.Where(sq.Eq{"refresh_token": req.RefreshToken})
	}

	query, args, err := q.ToSql()
	if err != nil {
		return fmt.Errorf("error building revoke grant sql: %w", err)
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("error excecuting revoke grant sql: %w", err)
	}

	return nil
}
