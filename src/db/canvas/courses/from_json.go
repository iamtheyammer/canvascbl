package courses

import (
	"encoding/json"
	"github.com/pkg/errors"
)

type CanvasCoursesResponse []struct {
	CourseCode string `json:"course_code"`
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	UUID       string `json:"uuid"`
	State      string `json:"workflow_state"`
}

func FromJSON(j *string) (*CanvasCoursesResponse, error) {
	var c CanvasCoursesResponse

	err := json.Unmarshal([]byte(*j), &c)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshaling into CanvasCoursesResponse")
	}

	return &c, nil
}
