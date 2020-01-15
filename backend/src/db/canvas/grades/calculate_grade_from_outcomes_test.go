package grades

import (
	"reflect"
	"testing"
)

func TestCalculateGradeFromOutcomeScores(t *testing.T) {
	type args struct {
		os []float64
	}
	tests := []struct {
		name string
		args args
		want Grade
	}{
		{
			"test_I",
			struct{ os []float64 }{os: []float64{1, 2, 3, 4, 5, 6, 7}},
			Grade{"I", 0, 0, 0},
		},
		{
			"test_B",
			struct{ os []float64 }{os: []float64{4, 4, 3.5, 3.3, 3, 1.8}},
			Grade{"B", 3, 2.6, 1.8},
		},
		{
			"test_B+",
			struct{ os []float64 }{os: []float64{4, 3.5, 2.6, 2.3, 2.2}},
			Grade{"B+", 4, 2.6, 2.2},
		},
		{
			"test_A",
			struct{ os []float64 }{os: []float64{4, 4, 3.5, 3.3, 3, 3}},
			Grade{"A", 6, 3.3, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalculateGradeFromOutcomeScores(tt.args.os); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CalculateGradeFromOutcomeScores() = %+v, want %+v", got.Grade, tt.want.Grade)
			}
		})
	}
}
