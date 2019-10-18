package grades

import (
	"github.com/Masterminds/squirrel"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
)

type InsertRequest struct {
	CourseID      string
	Grade         string
	UserLTIUserID string
}

func Insert(req InsertRequest) error {
	query, args, err := squirrel.Insert("grades").SetMap(map[string]interface{}{
		"course_id":        req.CourseID,
		"grade":            req.Grade,
		"user_lti_user_id": req.UserLTIUserID,
	}).ToSql()
	if err != nil {
		return errors.Wrap(err, "error building insert grade sql")
	}

	_, err = util.DB.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "error executing insert grade sql")
	}

	return nil
}
