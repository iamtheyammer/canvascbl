package db

import (
	"github.com/iamtheyammer/canvascbl/backend/src/db/canvas/outcomes"
	outcomessvc "github.com/iamtheyammer/canvascbl/backend/src/db/services/outcomes"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
)

// InsertOutcome inserts an outcome from a JSON outcome response from Canvas
func InsertOutcome(or *string) {
	cor, err := outcomes.OutcomeFromJSON(or)
	if err != nil {
		handleError(errors.Wrap(err, "error getting CanvasOutcomeResponse from JSON"))
		return
	}

	cID := uint64(0)
	if cor.ContextType == "Course" {
		cID = uint64(cor.ContextID)
	}

	err = outcomessvc.InsertOutcome(util.DB, &outcomessvc.InsertRequest{
		CanvasID:       uint64(cor.ID),
		CourseID:       cID,
		ContextID:      uint64(cor.ContextID),
		DisplayName:    cor.DisplayName,
		Title:          cor.Title,
		MasteryPoints:  cor.MasteryPoints,
		PointsPossible: cor.PointsPossible,
	})
	if err != nil {
		handleError(errors.Wrap(err, "error inserting an outcome"))
		return
	}

	return
}
