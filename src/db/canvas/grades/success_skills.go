package grades

import "strings"

func isSuccessSkillsOutcome(outcomeName string) bool {
	return strings.HasPrefix(outcomeName, "SS")
}
