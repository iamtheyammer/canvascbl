package gpas

import (
	"fmt"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
)

type InsertRequest struct {
	CanvasUserID     uint64
	Weighted         bool
	GPA              float64
	GPAWithSubgrades float64
	ManualFetch      bool
}

func InsertMultiple(db services.DB, req *[]InsertRequest) error {
	q := util.Sq.
		Insert("gpas").
		Columns(
			"canvas_user_id",
			"weighted",
			"gpa",
			"gpa_with_subgrades",
			"manual_fetch",
		)

	for _, g := range *req {
		q = q.Values(
			g.CanvasUserID,
			g.Weighted,
			g.GPA,
			g.GPAWithSubgrades,
			g.ManualFetch,
		)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return fmt.Errorf("error building insert multiple gpas sql: %w", err)
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("error executing insert multiple gpas sql: %w", err)
	}

	return nil
}
