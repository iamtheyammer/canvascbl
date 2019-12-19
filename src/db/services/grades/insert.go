package grades

import (
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
)

type InsertRequest struct {
	Grade string
	// int because -1 = same; 0 = no; 1 = yes
	HasSuccessSkills int
	CourseID         int
	UserCanvasID     int
}

func Insert(db services.DB, req *InsertRequest) error {
	q := util.Sq.
		Insert("grades")

	vals := map[string]interface{}{
		"course_id": req.CourseID,
		"grade":     req.Grade,
		// using an Expr because if I don't, it will set it to $1 instead of $3, which is required
		"user_canvas_id": req.UserCanvasID,
	}

	if req.HasSuccessSkills == 0 {
		vals["has_success_skills"] = false
	} else if req.HasSuccessSkills == 1 {
		vals["has_success_skills"] = true
	}

	q = q.SetMap(vals)

	query, args, err := q.ToSql()
	if err != nil {
		return errors.Wrap(err, "error building insert grade sql")
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "error executing insert grade sql")
	}

	return nil
}
