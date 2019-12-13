package grades

import (
	"encoding/json"
	"github.com/pkg/errors"
)

// CanvasRollupsResponse models the data returned from Canvas exactly. This should really only be used in relation
// to a Rollup.
type CanvasRollupsResponse struct {
	Rollups []struct {
		Links struct {
			User string `json:"user"`
		} `json:"links"`
		Scores []struct {
			Title string  `json:"title"`
			Score float64 `json:"score"`
			Links struct {
				Outcome string `json:"outcome"`
			}
		} `json:"scores"`
	} `json:"rollups"`
}

type OutcomeScore struct {
	Score           float64
	IsSuccessSkills bool
}

// GetCanvasRollupsResponseFromJsonString gets a CanvasRollupsResponse from a JSON string
func GetCanvasRollupsResponseFromJsonString(j *string) (*CanvasRollupsResponse, error) {
	var d CanvasRollupsResponse

	err := json.Unmarshal([]byte(*j), &d)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal string into CanvasRollupsResponse")
	}

	return &d, nil
}

// GetOutcomeScoresFromCanvasRollupsResponse gets an array of outcome scores from a CanvasRollupsResponse,
// often gotten from GetCanvasRollupsResponseFromJsonString.
func GetOutcomeScoresFromCanvasRollupsResponse(crr *CanvasRollupsResponse) (*[]OutcomeScore, error) {
	rs := crr.Rollups

	// there will be one rollup per user and we only do this by single user at the moment.
	if len(rs) != 1 {
		return nil, errors.New("error getting scores from canvas rollups: rollups.length != 1")
	}

	r := rs[0]

	var s []OutcomeScore

	for _, v := range r.Scores {
		s = append(s, OutcomeScore{Score: v.Score, IsSuccessSkills: isSuccessSkillsOutcome(v.Links.Outcome)})
	}

	return &s, nil
}
