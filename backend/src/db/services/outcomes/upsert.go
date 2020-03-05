package outcomes

import (
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
)

// UpsertOutcome upserts an outcome
func UpsertOutcome(db services.DB, req *InsertRequest) error {
	query, args, err := util.Sq.
		Insert("outcomes").
		SetMap(map[string]interface{}{
			"canvas_id":       req.CanvasID,
			"course_id":       req.CourseID,
			"context_id":      req.ContextID,
			"display_name":    req.DisplayName,
			"title":           req.Title,
			"mastery_points":  req.MasteryPoints,
			"points_possible": req.PointsPossible,
		}).
		Suffix("ON CONFLICT (canvas_id) DO UPDATE SET " +
			"display_name = EXCLUDED.display_name, " +
			"title = EXCLUDED.title, " +
			"mastery_points = EXCLUDED.mastery_points, " +
			"points_possible = EXCLUDED.points_possible",
		).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "error building upsert outcome sql")
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "error executing upsert outcome sql")
	}

	return nil
}
