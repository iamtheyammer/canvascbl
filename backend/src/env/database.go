package env

var (
	DatabaseDSN = getEnvOrPanic("DATABASE_DSN")
)
