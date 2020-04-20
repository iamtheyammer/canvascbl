package oauth2

type Scope string

const (
	ScopeProfile             = Scope("user_profile")
	ScopeObservees           = Scope("observees")
	ScopeCourses             = Scope("courses")
	ScopeAlignments          = Scope("alignments")
	ScopeAssignments         = Scope("assignments")
	ScopeOutcomeResults      = Scope("outcome_results")
	ScopeOutcomes            = Scope("outcomes")
	ScopeGrades              = Scope("grades")
	ScopeDetailedGrades      = Scope("detailed_grades")
	ScopePreviousGrades      = Scope("previous_grades")
	ScopeAverageGrade        = Scope("average_course_grade")
	ScopeAverageOutcomeScore = Scope("average_outcome_score")
	ScopeGPA                 = Scope("gpa")
	ScopeNotifications       = Scope("notifications")
	ScopeEnrollments         = Scope("enrollments")
)

// ValidateScopes ensures that all requested scopes are valid.
// It returns a bool (ok), and a string which will only be included if a scope
// is invalid, containing the invalid scope.
func ValidateScopes(scopes []string) (bool, *string) {
	// only one grades scope is allowed
	hasOneGrades := false

	for _, s := range scopes {
		switch Scope(s) {
		case ScopeProfile:
		case ScopeObservees:
		case ScopeCourses:
		case ScopeAlignments:
		case ScopeAssignments:
		case ScopeOutcomeResults:
		case ScopeOutcomes:
		case ScopeGrades:
			if hasOneGrades {
				return false, &s
			} else {
				hasOneGrades = true
			}
		case ScopeDetailedGrades:
			if hasOneGrades {
				return false, &s
			} else {
				hasOneGrades = true
			}
		case ScopePreviousGrades:
		case ScopeAverageGrade:
		case ScopeAverageOutcomeScore:
		case ScopeGPA:
		case ScopeNotifications:
		case ScopeEnrollments:
		default:
			return false, &s
		}
	}

	return true, nil
}
