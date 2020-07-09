package users

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
)

type UpdateUserRequest struct {
	WhereID           uint64
	WhereCanvasUserID uint64

	HasValidSubscription *bool
}

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

func UpdateUser(db services.DB, req *UpdateUserRequest) error {
	q := util.Sq.
		Update("users")

	if req.WhereID > 0 {
		q = q.Where(sq.Eq{"id": req.WhereID})
	}

	if req.WhereCanvasUserID > 0 {
		q = q.Where(sq.Eq{"canvas_user_id": req.WhereCanvasUserID})
	}

	if req.HasValidSubscription != nil {
		q = q.Set("has_valid_subscription", *req.HasValidSubscription)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return fmt.Errorf("error building update user sql: %w", err)
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("error executing update user sql: %w", err)
	}

	return nil
}
