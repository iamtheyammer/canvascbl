package util

import (
	"net/http"
)

const RequestDetailsFailedValidationMessage = "no canvas token or subdomain not allowed"

// RequestDetails contains details needed by functions to perform requests.
// While this type is exported, you shouldn't use it-- use NewRequestDetails to write to it.
type RequestDetails struct {
	Token string
}

// NewRequestDetails is the way to create a RequestDetails object for use with other
// canvasapi functions.
func NewRequestDetails(
	token string,
) RequestDetails {

	rd := RequestDetails{
		Token: token,
	}

	return rd
}

func GetRequestDetailsFromRequest(r *http.Request) (bool, *RequestDetails) {
	var token string

	// token - header
	headerToken := r.Header.Get("x-canvas-token")
	if len(headerToken) > 0 {
		token = headerToken
	} else {
		// in query
		keys, ok := r.URL.Query()["token"]

		if ok && len(keys[0]) > 1 {
			// returns an array of results
			token = keys[0]
		}
	}

	if len(token) < 1 {
		return false, nil
	}

	rd := NewRequestDetails(token)

	// not specified
	return true, &rd
}
