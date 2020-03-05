package courses

import (
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
)

type OutcomeRollupInsertRequest struct {
	CanvasUserID  uint64
	CourseID      uint64
	OutcomeID     uint64
	Score         float64
	TimesAssessed uint64
}

type OutcomeResultInsertRequest struct {
	ID              uint64
	CourseID        uint64
	AssignmentID    uint64
	OutcomeID       uint64
	UserID          uint64
	AchievedMastery bool
	Score           float64
	Possible        float64
	SubmissionTime  string
}

func InsertMultipleOutcomeRollups(db services.DB, req *[]OutcomeRollupInsertRequest) error {
	q := util.Sq.
		Insert("outcome_rollups").
		Columns(
			"course_canvas_id",
			"outcome_canvas_id",
			"user_canvas_id",
			"score",
			"times_assessed",
		)

	for _, or := range *req {
		timesAssessed := &or.TimesAssessed
		if or.TimesAssessed == 0 {
			timesAssessed = nil
		}

		q = q.Values(
			or.CourseID,
			or.OutcomeID,
			or.CanvasUserID,
			or.Score,
			timesAssessed,
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

func InsertMultipleOutcomeResults(db services.DB, req *[]OutcomeResultInsertRequest) error {
	q := util.Sq.
		Insert("outcome_results").
		Columns(
			"canvas_id",
			"course_canvas_id",
			"assignment_canvas_id",
			"outcome_canvas_id",
			"user_canvas_id",
			"achieved_mastery",
			"score",
			"possible",
			"submission_time",
		).
		Suffix(
			"ON CONFLICT ON CONSTRAINT outcome_results_canvas_id_key DO UPDATE SET " +
				"achieved_mastery = excluded.achieved_mastery, " +
				"score = excluded.score, " +
				"possible = excluded.possible, " +
				"submission_time = excluded.submission_time",
		)

	for _, or := range *req {
		q = q.Values(
			or.ID,
			or.CourseID,
			or.AssignmentID,
			or.OutcomeID,
			or.UserID,
			or.AchievedMastery,
			or.Score,
			or.Possible,
			or.SubmissionTime,
		)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return errors.Wrap(err, "error building insert multiple outcome results sql")
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "error executing insert multiple outcome results sql")
	}

	return nil
}
