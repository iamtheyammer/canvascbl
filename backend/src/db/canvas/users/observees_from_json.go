package users

import (
	"encoding/json"
	"github.com/pkg/errors"
)

type CanvasObserveesResponse []struct {
	CreatedAt                     string  `json:"created_at"`
	ID                            uint64  `json:"id"`
	Name                          string  `json:"name"`
	ObservationLinkRootAccountIds []int64 `json:"observation_link_root_account_ids"`
	ShortName                     string  `json:"short_name"`
	SortableName                  string  `json:"sortable_name"`
}

func ObserveesFromJSON(o *string) (*CanvasObserveesResponse, error) {
	var c CanvasObserveesResponse

	err := json.Unmarshal([]byte(*o), &c)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal observees from JSON")
	}

	return &c, nil
}
