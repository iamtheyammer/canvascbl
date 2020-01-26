package db

import (
	outcomessvc "github.com/iamtheyammer/canvascbl/backend/src/db/services/outcomes"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
	"time"
)

var memoizedOutcomeAverages = map[uint64]struct {
	OutcomeAverage *outcomessvc.OutcomeAverage
	ValidUntil     time.Time
}{}

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

func GetUserMostRecentOutcomeRollupScore(userCanvasIDs []uint64) (*outcomessvc.OutcomeRollupScore, error) {
	return outcomessvc.GetUserMostRecentScore(util.DB, userCanvasIDs)
}
