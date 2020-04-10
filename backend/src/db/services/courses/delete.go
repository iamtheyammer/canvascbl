package courses

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
)

func Show(db services.DB, req *HideRequest) error {
	query, args, err := util.Sq.
		Delete("hidden_courses").
		Where(sq.Eq{"user_id": req.UserID}).
		Where(sq.Eq{"course_id": req.CourseID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("error building show course sql: %w", err)
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("error executing show course sql: %w", err)
	}

	return nil
}
