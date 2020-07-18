package util

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func Healthcheck(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := DB.Ping()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("error pinging database"))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("200 OK"))
}
