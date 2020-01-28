package canvas_tokens

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
)

type DeleteRequest struct {
	Token        string
	RefreshToken string
}

func Delete(db services.DB, req *DeleteRequest) error {
	q := util.Sq.
		Delete("canvas_tokens")

	if len(req.Token) > 0 {
		q = q.Where(sq.Eq{"token": req.Token})
	}

	if len(req.RefreshToken) > 0 {
		q = q.Where(sq.Eq{"refresh_token": req.RefreshToken})
	}

	query, args, err := q.ToSql()
	if err != nil {
		return errors.Wrap(err, "error building delete canvas tokens sql")
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "error executing delete canvas tokens sql")
	}

	return nil
}
