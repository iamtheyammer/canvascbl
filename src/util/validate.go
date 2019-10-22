package util

import (
	"github.com/iamtheyammer/canvascbl/backend/src/env"
	"regexp"
	"strings"
)

var lowercaseStringRegex = regexp.MustCompile("[a-z]{1,64}")
var numberRegex = regexp.MustCompile("[0-9.,]{1,32}")
var uuidRegex = regexp.MustCompile("([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}){1}")

// ValidateSubdomain takes a subdomain and ensures it's allowed per the environment variable.
func ValidateSubdomain(req string) bool {
	allowed := env.ProxyAllowedSubdomains

	if allowed[0] == "*" {
		return true
	}

	for _, sd := range allowed {
		if sd == req {
			return true
		}
	}

	return false
}

// ValidateLowercaseString validates that the string contains 1 to 64 lowercase english letters.
func ValidateLowercaseString(req string) bool {
	// the req must be 1-64 lowercase letters
	return lowercaseStringRegex.FindString(req) == req
}

// ValidateIntegerString validates that the string contains 1 to 32 numbers (no commas, periods, etc.).
func ValidateIntegerString(req string) bool {
	return numberRegex.FindString(req) == req
}

// ValidateIncludes makes sure that the includes param doesn't include any of ?, & or /.
func ValidateIncludes(req string) bool {
	return !(strings.Contains(req, "?") ||
		strings.Contains(req, "&") ||
		strings.Contains(req, "/"))
}

func ValidateUUIDString(req string) bool {
	return uuidRegex.FindString(req) == req
}
