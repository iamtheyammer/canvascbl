package grades

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
)

type InsertRequest struct {
	Grade    string
	CourseID int
	UserID   int
}

func Insert(db *sql.DB, req *InsertRequest) error {
	query, args, err := util.Sq.
		Insert("grades").
		SetMap(map[string]interface{}{
			"course_id": req.CourseID,
			"grade":     req.Grade,
			// using an Expr because if I don't, it will set it to $1 instead of $3, which is required
			"user_lti_user_id": sq.Expr("(SELECT lti_user_id FROM users WHERE canvas_user_id=?)", req.UserID),
		}).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "error building insert grade sql")
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "error executing insert grade sql")
	}

	return nil
}
