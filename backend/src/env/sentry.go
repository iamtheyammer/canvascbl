package env

var (
	SentryDSN = getEnv("SENTRY_DSN", "https://no_op@no.op/no_op")
)
