package gradesapi

import (
	"reflect"
	"testing"
)

func Test_calculateGradeFromOutcomeResults(t *testing.T) {
	type args struct {
		results       map[uint64][]canvasOutcomeResult
		isAfterCutoff bool
	}
	tests := []struct {
		name string
		args args
		want *computedGrade
	}{
		{
			name: "test_4_4_3",
			args: struct {
				results       map[uint64][]canvasOutcomeResult
				isAfterCutoff bool
			}{results: map[uint64][]canvasOutcomeResult{
				1: {canvasOutcomeResult{Score: 4}},
				2: {canvasOutcomeResult{Score: 4}},
				3: {canvasOutcomeResult{Score: 3}}},
			}, want: &computedGrade{
				Grade: grade{"A", 6, 3.3, 3},
				Averages: map[uint64]computedAverage{
					1: {
						DidDropWorstScore: false,
						Average:           4,
					},
					2: {
						DidDropWorstScore: false,
						Average:           4,
					},
					3: {
						DidDropWorstScore: false,
						Average:           3,
					},
				},
			},
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateGradeFromOutcomeResults(tt.args.results, tt.args.isAfterCutoff); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("calculateGradeFromOutcomeResults() = %v, want %v", got, tt.want)
			}
		})
	}
}
