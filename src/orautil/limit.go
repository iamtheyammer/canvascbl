package orautil

import "fmt"

// BuildLimitExpression builds a limit expression for Oracle.
// It does not parameterize your limit, so it should not be from
// a user. The limit expression does not have extra spaces on either side.
func BuildLimitExpression(limit uint64) string {
	return fmt.Sprintf("FETCH NEXT %d ROWS", limit)
}
