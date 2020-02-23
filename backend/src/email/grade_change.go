package email

import (
	"strings"
)

// GradeChangeEmailData represents the data needed to send a grade change email.
type GradeChangeEmailData struct {
	To            string
	Name          string
	ClassName     string
	PreviousGrade string
	CurrentGrade  string
}

// ParentGradeChangeEmailData represents the data needed to send a grade change email to a parent.
type ParentGradeChangeEmailData struct {
	To string
	// Parent's name
	Name          string
	StudentName   string
	ClassName     string
	PreviousGrade string
	CurrentGrade  string
}

// SendGradeChangeEmail sends a grade change email to a student.
func SendGradeChangeEmail(req *GradeChangeEmailData) {
	firstName := strings.Split(req.Name, " ")[0]

	send(
		gradeChange,
		map[string]interface{}{
			"first_name":     firstName,
			"class_name":     req.ClassName,
			"previous_grade": req.PreviousGrade,
			"current_grade":  req.CurrentGrade,
		},
		req.To,
		req.Name,
	)
}

// SendParentGradeChangeEmail sends a grade change email to a parent.
func SendParentGradeChangeEmail(req *ParentGradeChangeEmailData) {
	firstName := strings.Split(req.Name, " ")[0]
	studentFirstName := strings.Split(req.StudentName, " ")[0]

	send(
		parentGradeChange,
		map[string]interface{}{
			"student_first_name": studentFirstName,
			"first_name":         firstName,
			"class_name":         req.ClassName,
			"previous_grade":     req.PreviousGrade,
			"current_grade":      req.CurrentGrade,
		},
		req.To,
		req.Name,
	)
}
