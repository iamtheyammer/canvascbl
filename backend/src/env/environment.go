package env

type Environment string

const (
	EnvironmentDevelopment Environment = "development"
	EnvironmentStaging     Environment = "staging"
	EnvironmentProduction  Environment = "production"
)

var (
	Env = Environment(getEnvOrPanic("ENVIRONMENT"))
)
