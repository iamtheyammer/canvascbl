package grades

import (
	"encoding/json"
	"github.com/pkg/errors"
	"strconv"
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
func GetOutcomeScoresFromCanvasRollupsResponse(crr *CanvasRollupsResponse) (*map[uint64][]float64, error) {
	rs := crr.Rollups

	ret := map[uint64][]float64{}

	for _, r := range rs {
		uID, err := strconv.Atoi(r.Links.User)
		if err != nil {
			return nil, errors.Wrap(err, "error converting user id into an int")
		}

		var s []float64
		for _, v := range r.Scores {
			s = append(s, v.Score)
		}

		ret[uint64(uID)] = s
	}

	return &ret, nil
}
