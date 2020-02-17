package env

import "fmt"

var (
	// HTTPPort is the port for the HTTP Server to run on.
	HTTPPort = fmt.Sprintf(":%s", getEnv("PORT", "8000"))
	// BaseURL is the API's URL, like "https://api.canvascbl.com"
	BaseURL = getEnvOrPanic("BASE_URL")
	// ScriptKey is used for autonomous tasks that need a key to run.
	ScriptKey = getEnv("SCRIPT_KEY", "default_script_key")
)
