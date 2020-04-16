package util

import "strings"

var scopes = []string{
	"url:GET|/api/v1/outcomes/:id",
	"url:GET|/api/v1/users/:id",
	"url:GET|/api/v1/users/:user_id/profile",
	"url:GET|/api/v1/courses",
	"url:GET|/api/v1/courses/:course_id/assignments",
	"url:GET|/api/v1/courses/:course_id/outcome_groups",
	"url:GET|/api/v1/courses/:course_id/outcome_groups/:id/outcomes",
	"url:GET|/api/v1/courses/:course_id/outcome_results",
	"url:GET|/api/v1/courses/:course_id/outcome_rollups",
	"url:GET|/api/v1/courses/:course_id/outcome_alignments",
	"url:GET|/api/v1/users/:user_id/observees",

	/* BEGIN BLOCK added 4/15/20 for teacher app */

	// courses
	"url:GET|/api/v1/users/:user_id/courses/:course_id/assignments",
	"url:GET|/api/v1/courses/:course_id/assignments/:id",
	"url:GET|/api/v1/users/:user_id/courses",
	"url:GET|/api/v1/courses/:course_id/effective_due_dates",
	"url:GET|/api/v1/courses/:course_id/users",
	"url:GET|/api/v1/courses/:course_id/search_users",
	"url:GET|/api/v1/courses/:course_id/recent_students",
	"url:GET|/api/v1/courses/:course_id/users/:id",
	"url:GET|/api/v1/courses/:course_id/todo",
	"url:GET|/api/v1/courses/:course_id/settings",

	// enrollments
	"url:GET|/api/v1/courses/:course_id/enrollments",
	"url:GET|/api/v1/sections/:section_id/enrollments",
	"url:GET|/api/v1/users/:user_id/enrollments",

	// grade change
	"url:GET|/api/v1/audit/grade_change/assignments/:assignment_id",
	"url:GET|/api/v1/audit/grade_change/courses/:course_id",
	"url:GET|/api/v1/audit/grade_change/students/:student_id",

	// gradebook history
	"url:GET|/api/v1/courses/:course_id/gradebook_history/days",

	// rubrics
	"url:GET|/api/v1/courses/:course_id/rubrics/:id",

	// submissions
	"url:GET|/api/v1/courses/:course_id/assignments/:assignment_id/submissions",
	"url:GET|/api/v1/courses/:course_id/quizzes/:quiz_id/submissions",
	"url:GET|/api/v1/courses/:course_id/students/submissions",
	"url:GET|/api/v1/courses/:course_id/assignments/:assignment_id/submissions/:user_id",
	"url:GET|/api/v1/courses/:course_id/assignments/:assignment_id/submission_summary",

	// outcome proficiency
	"url:GET|/api/v1/accounts/:account_id/outcome_proficiency",

	// analytics
	"url:GET|/api/v1/courses/:course_id/analytics/activity",
	"url:GET|/api/v1/courses/:course_id/analytics/assignments",
	"url:GET|/api/v1/courses/:course_id/analytics/student_summaries",
	"url:GET|/api/v1/courses/:course_id/analytics/users/:student_id/activity",
	"url:GET|/api/v1/courses/:course_id/analytics/users/:student_id/assignments",
	"url:GET|/api/v1/courses/:course_id/analytics/users/:student_id/communication",

	/* END BLOCK 4/15/20 */
}

var stringScopes = strings.Join(scopes, " ")

var googleScopes = []string{
	"https://www.googleapis.com/auth/userinfo.email",
	"https://www.googleapis.com/auth/userinfo.profile",
}

var stringGoogleScopes = strings.Join(googleScopes, " ")

// GetScopesList gets the static string list of scopes.
func GetScopesList() string {
	return stringScopes
}

// GetGoogleScopesList gets the static list of google scopes.
func GetGoogleScopesList() string {
	return stringGoogleScopes
}
