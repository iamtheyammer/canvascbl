package util

import (
	"fmt"
	"github.com/iamtheyammer/canvascbl/backend/src/env"
	"log"
	"net/http"
	"strings"
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

func HandleCanvasResponse(w http.ResponseWriter, resp *http.Response, body string) {
	sc := resp.StatusCode
	if sc < 200 || sc > 399 {
		SendCanvasError(w, resp, body)
		return
	}

	SendCanvasSuccess(w, resp, body)
}

func SendCanvasError(w http.ResponseWriter, resp *http.Response, efc string) {
	w.Header().Set("X-Canvas-Status-Code", fmt.Sprintf("%d", resp.StatusCode))

	if reqURL := resp.Request.URL.String(); !strings.Contains(reqURL, env.CanvasOAuth2ClientSecret) {
		w.Header().Set("X-Canvas-URL", reqURL)
	} else {
		w.Header().Set("X-Canvas-URL", "omitted")
	}

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	// can't use 502 with cloudflare
	//w.WriteHeader(http.StatusBadGateway)
	w.WriteHeader(CanvasProxyErrorCode)
	_, err := fmt.Fprint(w, efc)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func SendCanvasSuccess(w http.ResponseWriter, resp *http.Response, body string) {
	w.Header().Set("X-Canvas-Status-Code", fmt.Sprintf("%d", resp.StatusCode))

	if reqURL := resp.Request.URL.String(); !strings.Contains(reqURL, env.CanvasOAuth2ClientSecret) {
		w.Header().Set("X-Canvas-URL", reqURL)
	} else {
		w.Header().Set("X-Canvas-URL", "omitted")
	}

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.WriteHeader(http.StatusOK)
	_, err := fmt.Fprint(w, body)
	if err != nil {
		log.Fatal(err)
	}
	return
}

