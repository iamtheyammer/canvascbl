package email

import (
	"strings"
)

func SendGradeChangeEmail(
	to string,
	name string,
	className string,
	previousGrade string,
	currentGrade string,
) {
	names := strings.Split(name, " ")
	firstName := names[0]

	send(
		gradeChange,
		map[string]interface{}{
			"first_name":     firstName,
			"class_name":     className,
			"previous_grade": previousGrade,
			"current_grade":  currentGrade,
		},
		to,
		name,
	)
}
