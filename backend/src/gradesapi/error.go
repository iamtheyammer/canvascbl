package gradesapi

import (
	"encoding/json"
	"fmt"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"net/http"
)

// handleError handles an error and returns a response to the client.
// For Internal Server Errors, try handleISE(...)
func handleError(w http.ResponseWriter, gep gradesErrorResponse, httpStatusCode int) {
	resp, err := json.Marshal(&gep)
	if err != nil {
		util.HandleError(fmt.Errorf("error marshaling gradesErrorResponse: %w", err))
		util.SendInternalServerError(w)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	_, _ = w.Write(resp)
	return
}

// handleISE handles an Internal Server Error
func handleISE(w http.ResponseWriter, err error) {
	util.HandleError(fmt.Errorf("error in gradesHandler: %w", err))
	util.SendInternalServerError(w)
	return
}
