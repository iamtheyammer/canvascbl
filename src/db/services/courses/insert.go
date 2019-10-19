package courses

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
)

type AssignmentInsertRequest struct {
	CourseID uint64
	CanvasID uint64
	IsQuiz   bool
	Name     string
}

type OutcomeRollupInsertRequest struct {
	CourseID      uint64
	OutcomeID     uint64
	UserCanvasID  uint64
	Score         float64
	TimesAssessed uint64
}

func InsertMultipleAssignments(db *sql.DB, req *[]AssignmentInsertRequest) error {
	q := util.Sq.
		Insert("assignments").
		Columns(
			"course_id",
			"canvas_id",
			"is_quiz",
			"name",
		).
		Suffix("ON CONFLICT DO NOTHING")

	for _, a := range *req {
		q = q.Values(
			a.CourseID,
			a.CanvasID,
			a.IsQuiz,
			a.Name,
		)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return errors.Wrap(err, "error building insert multiple assignments sql")
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "error executing insert multiple assignments sql")
	}

	return nil
}

func InsertMultipleOutcomeRollups(db *sql.DB, uID uint64, req *[]OutcomeRollupInsertRequest) error {
	q := util.Sq.
		Insert("outcome_rollups").
		Prefix("WITH users_meta AS (SELECT lti_user_id FROM users WHERE canvas_user_id=?)", uID).
		Columns(
			"course_canvas_id",
			"outcome_canvas_id",
			"user_canvas_lti_id",
			"score",
			"times_assessed",
		)

	for _, or := range *req {
		q = q.Values(
			or.CourseID,
			or.OutcomeID,
			sq.Expr("(SELECT lti_user_id FROM users_meta)"),
			or.Score,
			or.TimesAssessed,
		)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return errors.Wrap(err, "error building insert multiple outcome rollups sql")
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "error executing insert multiple outcome rollups sql")
	}

	return nil
}
