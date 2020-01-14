package env

import "fmt"

var (
	HTTPPort = fmt.Sprintf(":%s", getEnv("PORT", "8000"))
)
