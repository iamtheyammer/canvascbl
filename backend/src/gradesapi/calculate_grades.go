package gradesapi

import (
	"math"
	"sort"
	"strings"
)

type grade struct {
	// the grade, like A or B
	Grade string `json:"grade"`
	// how "good" the grade is, higher is better
	Rank int8 `json:"rank"`
	// 75% = most; aka max
	MostAbove float64 `json:"most_above"`
	// aka min
	AllAbove float64 `json:"all_above"`
	// GPAVal represents the gpa points this grade is worth, for example, a B-, B, or B+ would be 3.0.
	// Not included with JSON
	GPAVal float64 `json:"-"`
	// SubgradeGPAVal represents the gpa points this grade is worth, traditionally, so an A- would be 3.7.
	// Not included with JSON
	SubgradeGPAVal float64 `json:"-"`
}

type computedAverage struct {
	// Whether the worst score was dropped for a better average
	DidDropWorstScore bool `json:"did_drop_worst_score"`
	// The computed final average
	Average float64 `json:"average"`
}

type computedGrade struct {
	// The grade for the course
	Grade grade `json:"grade"`
	// Relationship between outcome ID and whether the last score was dropped
	Averages map[uint64]computedAverage `json:"averages"`
}

type distanceLearningPair struct {
	CourseName               string `json:"course_name"`
	OriginalCourseID         uint64 `json:"original_course_id"`
	DistanceLearningCourseID uint64 `json:"distance_learning_course_id"`
}

// gradeMap is a slice of possible grades
var gradeMap = []grade{
	{"A", 6, 3.3, 3, 4, 4},
	{"A-", 5, 3.3, 2.5, 4, 3.7},
	{"B+", 4, 2.6, 2.2, 3, 3.3},
	{"B", 3, 2.6, 1.8, 3, 3},
	{"B-", 2, 2.6, 1.5, 3, 2.7},
	{"C", 1, 2.2, 1.5, 2, 2},
	{"I", 0, 0, 0, 0, 0},
}

var naGrade = grade{"N/A", -1, 0, 0, 0, 0}

/*
calculateGradeFromOutcomeResults calculates a grade object from a map of scores.
The map should look like this: map[outcomeID<uint64>][]scores<float64>.

isAfterCutoff represents whether the lowest score should be dropped to improve a grade.

The returned map displays the relationship between the outcome ID and whether the last score was dropped.

Note that this is a very expensive function: ~O(n^5+n), so consider use inside of a goroutine.
*/
func calculateGradeFromOutcomeResults(results map[uint64][]canvasOutcomeResult, isAfterCutoff bool) *computedGrade {
	if len(results) < 1 {
		return &computedGrade{
			Grade:    naGrade,
			Averages: nil,
		}
	}

	// all data
	averages := map[uint64]computedAverage{}
	// just floats (for grade calculation)
	var avgs []float64

	for oID, rs := range results {
		if len(rs) < 1 {
			continue
		}

		if len(rs) == 1 && rs[0].Score > 0 {
			// why do the work if there's only one score?
			averages[oID] = computedAverage{
				DidDropWorstScore: false,
				Average:           rs[0].Score,
			}

			avgs = append(avgs, rs[0].Score)
			continue
		}

		var scores []float64
		for _, s := range rs {
			if s.Score > 0 {
				scores = append(scores, s.Score)
			}
		}

		if len(scores) < 1 {
			continue
		}

		sortedScores := sort.Float64Slice(scores)
		sort.Sort(sort.Reverse(sortedScores))

		var total float64
		for _, s := range sortedScores {
			total += s
		}

		numScores := float64(sortedScores.Len())

		didDrop := false

		// average of all sortedScores
		allScoreAvg := total / numScores
		// average of all sortedScores except for lowest (last, as it's sorted)
		// zero by default because this is only calculated if there is more than 1 score
		noLastAvg := float64(0)
		if numScores > 1 {
			// total - lastScore / numberOfScores - 1 (so that the average is right)
			noLastAvg = (total - sortedScores[int(numScores-1)]) / (numScores - 1)
		}

		// default the average to using all sortedScores
		avg := allScoreAvg

		// if the "dropped" average helped the score AND it isn't after the cutoff, use it.
		if (noLastAvg > allScoreAvg) && !isAfterCutoff {
			avg = noLastAvg
			didDrop = true
		}

		averages[oID] = computedAverage{
			DidDropWorstScore: didDrop,
			Average:           avg,
		}
		avgs = append(avgs, avg)
	}

	if len(avgs) < 1 {
		return &computedGrade{
			Grade:    naGrade,
			Averages: nil,
		}
	}

	// what is 75% of len(s)
	outcomesOverMinNeeded := int(math.Floor(float64(75*len(avgs)) / float64(100)))

	// float64 outcome results
	sortedOutcomes := sort.Float64Slice(avgs)
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

		// student qualifies for this grade
		if v.Rank > finalGrade.Rank {
			finalGrade = v
		}
	}

	return &computedGrade{
		Grade:    finalGrade,
		Averages: averages,
	}
}

/*

calculateGPAFromDetailedGrades returns a gpa object from
a detailedGrades object.
*/
func calculateGPAFromDetailedGrades(g detailedGrades) gpa {
	finalGPA := gpa{}
	for uID, cs := range g {
		var (
			cGPA calculatedGPA
			//subgrade sum
			subSum float64
			// default sum
			defSum float64
			// courses with a valid grade
			validClasses float64
		)

		for _, c := range cs {
			if c.Grade != naGrade {
				validClasses += 1
			}
			subSum += c.Grade.SubgradeGPAVal
			defSum += c.Grade.GPAVal
		}

		cGPA.Unweighted.Subgrades = subSum / validClasses
		cGPA.Unweighted.Default = defSum / validClasses

		finalGPA[uID] = cGPA
	}

	return finalGPA
}

func findDistanceLearningCoursePairs(courses []canvasCourse) *[]distanceLearningPair {
	// find regular vs. distance learning courses
	distanceCourses := map[uint64]canvasCourse{}
	originalCourses := map[uint64]canvasCourse{}

	courseMatches := map[string][]uint64{}

	for _, c := range courses {
		// nickname? bye.
		splits := strings.Split(c.Name, " - ")
		if len(splits) != 3 {
			continue
		}

		// like "Biology" from "Biology - S2-DL - Teacher"
		courseTitle := splits[0]

		// add course ID
		courseMatches[courseTitle] = append(courseMatches[courseTitle], c.ID)

		if c.EnrollmentTermID == spring20DLEnrollmentTermID {
			distanceCourses[c.ID] = c
			continue
		} else {
			originalCourses[c.ID] = c
			continue
		}
	}

	var coursePairs []distanceLearningPair

	// filter for 2s that have one distance and one original
	for title, ids := range courseMatches {
		if len(ids) != 2 {
			continue
		}

		foundDistanceLearning := false
		foundOriginal := false

		pair := distanceLearningPair{CourseName: title}

		for _, id := range ids {
			if _, ok := distanceCourses[id]; ok {
				foundDistanceLearning = true
				pair.DistanceLearningCourseID = id
				continue
			} else if _, ok := originalCourses[id]; ok {
				foundOriginal = true
				pair.OriginalCourseID = id
				continue
			}
		}

		if !foundDistanceLearning || !foundOriginal {
			continue
		}

		coursePairs = append(coursePairs, pair)
	}

	return &coursePairs
}

/*
calculateDistanceLearningGrades calculates distance learning grades for a single user.

It returns a slice of distance learning grades for the user.
*/
func calculateDistanceLearningGrades(courses []canvasCourse, grades map[uint64]computedGrade) []distanceLearningGrade {
	// find regular vs. distance learning courses
	distanceCourses := map[uint64]canvasCourse{}
	originalCourses := map[uint64]canvasCourse{}

	courseMatches := map[string][]uint64{}

	for _, c := range courses {
		// nickname? bye.
		splits := strings.Split(c.Name, " - ")
		if len(splits) != 3 {
			continue
		}

		// like "Biology" from "Biology - S2-DL - Teacher"
		courseTitle := splits[0]

		// add course ID
		courseMatches[courseTitle] = append(courseMatches[courseTitle], c.ID)

		if c.EnrollmentTermID == spring20DLEnrollmentTermID {
			distanceCourses[c.ID] = c
			continue
		} else {
			originalCourses[c.ID] = c
			continue
		}
	}

	var coursePairs []distanceLearningPair

	// filter for 2s that have one distance and one original
	for title, ids := range courseMatches {
		if len(ids) != 2 {
			continue
		}

		foundDistanceLearning := false
		foundOriginal := false

		pair := distanceLearningPair{CourseName: title}

		for _, id := range ids {
			// make sure user has grade in course
			if _, ok := grades[id]; !ok {
				continue
			}

			if _, ok := distanceCourses[id]; ok {
				foundDistanceLearning = true
				pair.DistanceLearningCourseID = id
				continue
			} else if _, ok := originalCourses[id]; ok {
				foundOriginal = true
				pair.OriginalCourseID = id
				continue
			}
		}

		if !foundDistanceLearning || !foundOriginal {
			continue
		}

		coursePairs = append(coursePairs, pair)
	}

	// finally, compute grades
	var dlGrades []distanceLearningGrade

	for _, pair := range coursePairs {
		grade := distanceLearningGrade{
			CourseName: pair.CourseName,
			Grade: struct {
				Grade string `json:"grade"`
				Rank  int    `json:"rank"`
			}{},
			OriginalCourseID:         pair.OriginalCourseID,
			DistanceLearningCourseID: pair.DistanceLearningCourseID,
		}

		dlCourseGrade := grades[pair.DistanceLearningCourseID]
		originalCourseGrade := grades[pair.OriginalCourseID]

		if dlCourseGrade.Grade == naGrade || originalCourseGrade.Grade == naGrade {
			grade.Grade.Grade = "N/A"
			grade.Grade.Rank = -1

			dlGrades = append(dlGrades, grade)
			continue
		}

		dlGrade := dlCourseGrade.Grade.Grade
		originalGrade := originalCourseGrade.Grade
		pass := false

		// dl == I ? instant fail
		if dlGrade == "I" {
			pass = false
			// dl != I, original != I OK
		} else if dlGrade != "I" && originalGrade.Rank > 0 {
			pass = true
			// dl != I, and original == I...
		} else if dlGrade != "I" && originalGrade.Rank == 0 {
			if dlGrade != "I" && dlGrade != "C" {
				pass = true
			} else {
				pass = false
			}
		}

		if pass {
			grade.Grade.Grade = "P"
			grade.Grade.Rank = 1
		} else {
			grade.Grade.Grade = "I"
			grade.Grade.Rank = 0
		}

		dlGrades = append(dlGrades, grade)
	}

	return dlGrades
}
