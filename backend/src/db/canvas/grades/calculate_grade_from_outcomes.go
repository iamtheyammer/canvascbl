package grades

import (
	"math"
	"sort"
)

type Grade struct {
	// the Grade, like A or B
	Grade string
	// how "good" the Grade is, higher is better
	Rank uint8
	// 75% = most; aka max
	MostAbove float64
	// aka min
	AllAbove float64
}

// gradeMap is a slice of possible grades
var gradeMap = []Grade{
	{"A", 6, 3.3, 3},
	{"A-", 5, 3.3, 2.5},
	{"B+", 4, 2.6, 2.2},
	{"B", 3, 2.6, 1.8},
	{"B-", 2, 2.6, 1.5},
	{"C", 1, 2.2, 1.5},
	{"I", 0, 0, 0},
}

// CalculateGradeFromOutcomeScores calculates a Grade from outcome scores (os)
func CalculateGradeFromOutcomeScores(os []float64) Grade {
	// what is 75% of len(s)
	outcomesOverMinNeeded := int(math.Floor(float64(75*len(os)) / float64(100)))

	// float64 outcome scores
	sortedOutcomes := sort.Float64Slice(os)
	// sorts in place
	sort.Sort(sort.Reverse(sortedOutcomes))
	// sortedOutcomes are now sorted

	// first 75% of outcomes
	countedOutcomes := sortedOutcomes[:outcomesOverMinNeeded]

	// if there is only one graded outcome in a class, that outcome is counted
	// this also fixes an array[-1] bug
	if len(countedOutcomes) < 1 {
		countedOutcomes = sortedOutcomes
	}

	lowestCountedOutcome := countedOutcomes[len(countedOutcomes)-1]

	// overall
	lowestOutcome := sortedOutcomes[len(sortedOutcomes)-1]

	// default to an I
	finalGrade := gradeMap[6]

	for _, v := range gradeMap {
		// lowest outcome is over minimum (AllAbove)
		if v.AllAbove > lowestOutcome {
			continue
		}

		// lowestCountedOutcome must be above MostAbove
		if v.MostAbove > lowestCountedOutcome {
			continue
		}

		// student qualifies for this Grade
		if v.Rank > finalGrade.Rank {
			finalGrade = v
		}
	}

	return finalGrade
}
