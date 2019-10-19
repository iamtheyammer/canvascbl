package grades

import (
	"math"
	"sort"
)

type grade struct {
	// the grade, like A or B
	Grade string
	// how "good" the grade is, higher is better
	Rank uint8
	// 75% = most; aka max
	MostAbove float64
	// aka min
	AllAbove float64
}

// gradeMap is a map of possible grades
var gradeMap = []grade{
	{"A", 6, 3.5, 3},
	{"A-", 5, 3.5, 2.5},
	{"B+", 4, 3, 2.5},
	{"B", 3, 3, 2.25},
	{"B-", 2, 3, 2},
	{"C", 1, 2.5, 2},
	{"I", 0, 0, 0},
}

// CalculateGradeFromOutcomeScores calculates a grade from outcome scores (os)
func CalculateGradeFromOutcomeScores(os []float64) string {
	// what is 75% of len(s)
	outcomesOverMinNeeded := int(math.Ceil(float64(75*len(os)) / float64(100)))

	// float64 outcome scores
	sortedOutcomes := sort.Float64Slice(os)
	// sorts in place
	sort.Sort(sort.Reverse(sortedOutcomes))
	// sortedOutcomes are now sorted

	// first 75% of outcomes
	countedOutcomes := sortedOutcomes[:outcomesOverMinNeeded]

	lowestCountedOutcome := countedOutcomes[len(countedOutcomes)-1]

	// overall
	lowestOutcome := sortedOutcomes[len(sortedOutcomes)-1]

	// in Golang, maps are not ordered
	// this means that we need to get all qualifying grades and sort.
	var finalGrade grade

	for _, v := range gradeMap {
		// lowest outcome is over minimum (AllAbove)
		if v.AllAbove > lowestOutcome {
			continue
		}

		// lowestCountedOutcome must be above MostAbove
		if v.MostAbove > lowestCountedOutcome {
			continue
		}

		// student qualifies for this grade
		if v.Rank > finalGrade.Rank {
			finalGrade = v
		}
	}

	return finalGrade.Grade
}
