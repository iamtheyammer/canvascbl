package main

import (
	"fmt"
	"os"
)

var (
	apiURLEnv   = getEnv("API_URL", "https://api.canvascbl.com/api/v1/grades/fetch_all")
	scriptKey   = getEnvOrPanic("SCRIPT_KEY")
	environment = getEnv("ENVIRONMENT", "unknown")
	awsRegion   = getEnv("AWS_REGION", "us-east-2")
)

func getEnvOrPanic(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Sprintf("Missing required environment variable '%v'\n", key))
	}
	return value
}

func getEnv(key string, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}

	return fallback
}
