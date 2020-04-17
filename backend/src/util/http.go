package util

import (
	"fmt"
	"log"
	"net/http"
)

const CanvasProxyErrorCode = 450

func SendJSONResponse(w http.ResponseWriter, j []byte) {
	w.Header().Set("Content-Type", "application/json")
	_, err := fmt.Fprint(w, string(j))
	if err != nil {
		log.Fatal(err)
	}
	return
}
