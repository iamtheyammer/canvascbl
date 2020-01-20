package canvas_tokens

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"time"
)

func UpdateFromRefreshToken(db services.DB, tokenID uint64, newToken string, expiresAt *time.Time) error {
	query, args, err := util.Sq.
		Update("canvas_tokens").
		Set("token", newToken).
		Set("expires_at", expiresAt).
		Where(sq.Eq{"id": tokenID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("error building update canvas token from refresh token sql: %w", err)
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("error executing update canvas token from refresh token sql: %w", err)
	}

	return nil
}
