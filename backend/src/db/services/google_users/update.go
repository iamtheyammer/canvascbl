package google_users

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
)

type UpdateRequest struct {
	UsersID uint64

	WhereID    uint64
	WhereEmail string
}

func Update(db services.DB, req *UpdateRequest) error {
	q := util.Sq.
		Update("google_users")

	if req.WhereID > 0 {
		q = q.Where(sq.Eq{"id": req.WhereID})
	}

	if len(req.WhereEmail) > 0 {
		q = q.Where(sq.Eq{"email": req.WhereEmail})
	}

	if req.UsersID > 0 {
		q = q.Set("users_id", req.UsersID)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return errors.Wrap(err, "error building update google user sql")
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "error executing update google user sql")
	}

	return nil
}
