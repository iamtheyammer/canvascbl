package grades

import "strings"

func isSuccessSkillsOutcome(title string) bool {
	t := strings.ToLower(title)

	return strings.HasPrefix(t, "SS") ||
		(strings.Contains(t, "teacher") && strings.Contains(t, "evaluation")) ||
		(strings.Contains(t, "teacher") && strings.Contains(t, "assessment")) ||
		(strings.Contains(t, "self") && strings.Contains(t, "evaluation")) ||
		(strings.Contains(t, "self") && strings.Contains(t, "assessment")) ||
		(strings.Contains(t, "peer") && strings.Contains(t, "evaluation")) ||
		(strings.Contains(t, "peer") && strings.Contains(t, "assessment"))
}
