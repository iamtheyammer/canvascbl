package outcomes

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
	"time"
)

type OutcomeAverage struct {
	OutcomeCanvasID uint64
	AverageScore    float64
	NumFactors      uint64
}

type OutcomeRollupScore struct {
	ID              uint64
	CourseCanvasID  uint64
	OutcomeCanvasID uint64
	UserCanvasID    uint64
	Score           float64
	TimesAssessed   uint64
	InsertedAt      time.Time
}

func GetUserMostRecentScore(db services.DB, userCanvasID uint64) (*OutcomeRollupScore, error) {
	query, args, err := util.Sq.
		Select(
			"id",
			"course_canvas_id",
			"outcome_canvas_id",
			"user_canvas_id",
			"score",
			"times_assessed",
			"inserted_at",
		).
		From("outcome_rollups").
		Where(sq.Eq{"user_canvas_id": userCanvasID}).
		OrderBy("inserted_at DESC").
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "error building get user most recent outcome rollup score sql")
	}

	row := db.QueryRow(query, args...)

	var ors OutcomeRollupScore
	err = row.Scan(
		&ors.ID,
		&ors.CourseCanvasID,
		&ors.OutcomeCanvasID,
		&ors.UserCanvasID,
		&ors.Score,
		&ors.TimesAssessed,
		&ors.InsertedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "error executing get user most recent outcome rollup score sql")
	}

	return &ors, nil
}

func GetAverage(db services.DB, outcomeID uint64) (*OutcomeAverage, error) {
	query, args, err := util.Sq.
		Select(
			"ROUND(AVG(average_score), 2) AS average_score",
			"COUNT(*) AS num_factors",
		).
		Prefix(`WITH outcomes_meta AS (
	SELECT DISTINCT ON (users.canvas_user_id)
		outcome_canvas_id AS outcome_id,
		AVG(score) AS average_score,
		COUNT(*) AS num_factors
	FROM
		outcome_rollups
		JOIN users ON users.canvas_user_id = outcome_rollups.user_canvas_id
	WHERE
		outcome_canvas_id = ?
		AND outcome_rollups.inserted_at > NOW() - interval '24 hours'
	GROUP BY
		outcome_canvas_id,
		users.canvas_user_id
)`, outcomeID).
		From("outcomes_meta").
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "error building get outcome average sql")
	}

	row := db.QueryRow(query, args...)

	var oa OutcomeAverage
	oa.OutcomeCanvasID = outcomeID
	err = row.Scan(&oa.AverageScore, &oa.NumFactors)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrap(err, "error executing get outcome average sql")
	}

	return &oa, nil
}
