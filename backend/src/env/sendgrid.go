package env

var (
	SendGridAPIKey = getEnvOrPanic("SENDGRID_API_KEY")
)
