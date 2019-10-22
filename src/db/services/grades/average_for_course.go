package grades

import (
	"database/sql"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
)

// CourseGradeAverage represents the average grade for a course, represented as a float.
type CourseGradeAverage struct {
	// NumInputs is the number of students' grades factored into the average
	NumInputs uint64
	// Average is the average grade in the class, as a float. 0 represents all I's and 6 represents all A's.
	// See src/db/canvas/grades/calculate_grade_from_outcomes.go for the rank of each grade.
	Average float64
}

func GetAverageForCourse(db services.DB, courseID uint64) (*CourseGradeAverage, error) {
	query, args, err := util.Sq.
		Select("COUNT(*) AS num_inputs", "AVG(grade_ints.grade_int) AS avg").
		Prefix(
			"WITH grade_ints AS(SELECT DISTINCT ON(user_lti_user_id)(CASE WHEN grade='I' "+
				"THEN 0 WHEN grade='C' THEN 1 WHEN grade='B-' THEN 2 WHEN grade='B' THEN 3 "+
				"WHEN grade='B+' THEN 4 WHEN grade='A-' THEN 5 WHEN grade='A' THEN 6 ELSE 0 END)"+
				"AS grade_int FROM grades WHERE inserted_at>NOW()-interval'24 hours' AND course_id=? "+
				"GROUP BY grades.grade,grades.user_lti_user_id,grades.inserted_at ORDER BY "+
				"user_lti_user_id,inserted_at DESC)",
			courseID).
		From("grade_ints").
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "error building get average for course sql")
	}

	row := db.QueryRow(query, args...)

	var (
		ga  CourseGradeAverage
		avg sql.NullFloat64
	)
	err = row.Scan(&ga.NumInputs, &avg)
	if err != nil {
		return nil, errors.Wrap(err, "error scanning course grade average")
	}

	if avg.Valid {
		ga.Average = avg.Float64
	}

	return &ga, nil
}
