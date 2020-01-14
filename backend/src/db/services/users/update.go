package users

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
)

func SoftDeleteUserObservees(db services.DB, observeeUserIDs []uint64) error {
	query, args, err := util.Sq.
		Update("observees").
		Where(sq.Eq{"observee_canvas_user_id": observeeUserIDs}).
		Set("deleted_at", sq.Expr("NOW()")).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "error building delete observees sql")
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "error executing delete observees sql")
	}

	return nil
}

func UnSoftDeleteUserObservees(db services.DB, observeeUserIDs []uint64) error {
	query, args, err := util.Sq.
		Update("observees").
		Where(sq.Eq{"observee_canvas_user_id": observeeUserIDs}).
		Set("deleted_at", nil).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "error building un-soft delete user observees sql")
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "error executing un-soft delete user observees sql")
	}

	return nil
}
