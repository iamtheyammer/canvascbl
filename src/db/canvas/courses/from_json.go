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

type CanvasAssignmentsResponse []struct {
	CourseID           uint64 `json:"course_id"`
	GradingType        string `json:"grading_type"`
	ID                 uint64 `json:"id"`
	IsQuizAssignment   bool   `json:"is_quiz_assignment"`
	Name               string `json:"name"`
	OmitFromFinalGrade bool   `json:"omit_from_final_grade"`
}

type CanvasOutcomeRollupsResponse struct {
	Rollups []struct {
		Links struct {
			User string `json:"user"`
		} `json:"links"`
		Scores []struct {
			Count int64 `json:"count"`
			Links struct {
				Outcome string `json:"outcome"`
			} `json:"links"`
			Score float64 `json:"score"`
		} `json:"scores"`
	} `json:"rollups"`
}

func FromJSON(j *string) (*CanvasCoursesResponse, error) {
	var c CanvasCoursesResponse

	err := json.Unmarshal([]byte(*j), &c)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshaling into CanvasCoursesResponse")
	}

	return &c, nil
}

func AssignmentsFromJSON(j *string) (*CanvasAssignmentsResponse, error) {
	var a CanvasAssignmentsResponse

	err := json.Unmarshal([]byte(*j), &a)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshaling into CanvasAssignmentsResponse")
	}

	return &a, nil
}

func OutcomeRollupsFromJSON(j *string) (*CanvasOutcomeRollupsResponse, error) {
	var or CanvasOutcomeRollupsResponse

	err := json.Unmarshal([]byte(*j), &or)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshaling into CanvasOutcomeRollupsResponse")
	}

	return &or, nil
}
