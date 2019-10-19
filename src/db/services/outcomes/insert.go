package outcomes

import (
	"database/sql"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
)

type InsertRequest struct {
	CanvasID       uint64
	CourseID       uint64
	ContextID      uint64
	DisplayName    string
	Title          string
	MasteryPoints  float64
	PointsPossible float64
}

// InsertOutcome inserts an outcome
func InsertOutcome(db *sql.DB, req *InsertRequest) error {
	vals := map[string]interface{}{
		"canvas_id":       req.CanvasID,
		"course_id":       req.CourseID,
		"context_id":      req.ContextID,
		"display_name":    req.DisplayName,
		"title":           req.Title,
		"mastery_points":  req.MasteryPoints,
		"points_possible": req.PointsPossible,
	}

	if req.CourseID == 0 {
		vals["course_id"] = nil
	}

	query, args, err := util.Sq.
		Insert("outcomes").
		SetMap(vals).
		Suffix("ON CONFLICT DO NOTHING").
		ToSql()

	if err != nil {
		return errors.Wrap(err, "error building insert outcome sql")
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "error executing insert outcome sql")
	}

	return nil
}
