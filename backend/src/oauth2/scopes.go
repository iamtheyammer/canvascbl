package oauth2

type Scope string

const (
	ScopeProfile             = Scope("profile")
	ScopeObservees           = Scope("observees")
	ScopeCourses             = Scope("courses")
	ScopeAlignments          = Scope("alignments")
	ScopeRollups             = Scope("rollups")
	ScopeAssignments         = Scope("assignments")
	ScopeOutcomes            = Scope("outcomes")
	ScopeGrades              = Scope("grades")
	ScopePreviousGrades      = Scope("previous_grades")
	ScopeAverageGrade        = Scope("average_course_grade")
	ScopeAverageOutcomeScore = Scope("average_outcome_score")
)

// ValidateScopes ensures that all requested scopes are valid.
// It returns a bool (ok), and a string which will only be included if a scope
// is invalid, containing the invalid scope.
func ValidateScopes(scopes []string) (bool, *string) {
	for _, s := range scopes {
		switch Scope(s) {
		case ScopeProfile:
		case ScopeObservees:
		case ScopeCourses:
		case ScopeAlignments:
		case ScopeRollups:
		case ScopeAssignments:
		case ScopeOutcomes:
		case ScopeGrades:
		case ScopePreviousGrades:
		case ScopeAverageGrade:
		case ScopeAverageOutcomeScore:
		default:
			return false, &s
		}
	}

	return true, nil
}
