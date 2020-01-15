package util

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type errorJsonResponse struct {
	Error string `json:"error"`
}

func SendUnauthorized(w http.ResponseWriter, reason string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	jError, _ := json.Marshal(&errorJsonResponse{Error: reason})
	SendJSONResponse(w, jError)
	return
}

func SendBadRequest(w http.ResponseWriter, reason string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	jError, _ := json.Marshal(&errorJsonResponse{Error: reason})
	SendJSONResponse(w, jError)
	return
}

func SendInternalServerError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	jError, _ := json.Marshal(&errorJsonResponse{Error: "Internal Server Error"})
	SendJSONResponse(w, jError)
	return
}

func SendMethodNotAllowed(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusMethodNotAllowed)
	jError, _ := json.Marshal(&errorJsonResponse{Error: "Method Not Allowed"})
	SendJSONResponse(w, jError)
	return
}

func SendNotFound(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	jError, _ := json.Marshal(&errorJsonResponse{Error: "Not Found"})
	SendJSONResponse(w, jError)
	return
}

func SendNotFoundWithReason(w http.ResponseWriter, reason string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	jError, _ := json.Marshal(&errorJsonResponse{Error: reason})
	SendJSONResponse(w, jError)
	return
}

func SendNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
	return
}

func SendRedirect(w http.ResponseWriter, to string) {
	w.Header().Set("Location", to)
	w.WriteHeader(http.StatusFound)

	_, err := fmt.Fprint(w, to)
	if err != nil {
		log.Fatal(err)
	}
	return
}
