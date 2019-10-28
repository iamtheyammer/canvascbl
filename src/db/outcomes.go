package db

import (
	"github.com/iamtheyammer/canvascbl/backend/src/db/canvas/outcomes"
	outcomessvc "github.com/iamtheyammer/canvascbl/backend/src/db/services/outcomes"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
	"time"
)

var memoizedOutcomeAverages = map[uint64]struct {
	OutcomeAverage *outcomessvc.OutcomeAverage
	ValidUntil     time.Time
}{}

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

func GetMemoizedOutcomeAverage(outcomeID uint64) (*outcomessvc.OutcomeAverage, error) {
	v, ok := memoizedOutcomeAverages[outcomeID]
	if ok && v.ValidUntil.After(time.Now()) {
		if v.OutcomeAverage == nil {
			return nil, nil
		}

		return v.OutcomeAverage, nil
	}

	// either there is no stored average or it has expired
	oa, err := outcomessvc.GetAverage(util.DB, outcomeID)
	if err != nil {
		return nil, errors.Wrap(err, "error getting outcome average")
	}

	memoizedOutcomeAverages[outcomeID] = struct {
		OutcomeAverage *outcomessvc.OutcomeAverage
		ValidUntil     time.Time
	}{OutcomeAverage: oa, ValidUntil: time.Now().Add(time.Minute * 5)}

	return oa, nil
}

func GetUserMostRecentOutcomeRollupScore(userLTIUserID string) (*outcomessvc.OutcomeRollupScore, error) {
	return outcomessvc.GetUserMostRecentScore(util.DB, userLTIUserID)
}
