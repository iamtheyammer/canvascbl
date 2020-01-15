package outcomes

import (
	"encoding/json"
	"github.com/pkg/errors"
)

type CanvasOutcomeResponse struct {
	ContextID      int64   `json:"context_id"`
	ContextType    string  `json:"context_type"`
	DisplayName    string  `json:"display_name"`
	ID             int64   `json:"id"`
	MasteryPoints  float64 `json:"mastery_points"`
	PointsPossible float64 `json:"points_possible"`
	Title          string  `json:"title"`
}

func OutcomeFromJSON(oj *string) (*CanvasOutcomeResponse, error) {
	var o CanvasOutcomeResponse

	err := json.Unmarshal([]byte(*oj), &o)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshaling outcome JSON")
	}

	return &o, nil
}
