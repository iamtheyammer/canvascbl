package canvas_tokens

import (
	"github.com/Masterminds/squirrel"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
)

type DeleteRequest struct {
	Token string
}

func Delete(db services.DB, req *DeleteRequest) error {
	query, args, err := util.Sq.
		Delete("canvas_tokens").
		Where(squirrel.Eq{"token": req.Token}).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "error building delete canvas tokens sql")
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "error executing delete canvas tokens sql")
	}

	return nil
}
