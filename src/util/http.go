package util

import (
	"fmt"
	"github.com/iamtheyammer/canvascbl/backend/src/env"
	"log"
	"net/http"
	"net/url"
	"strings"
)

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
	w.WriteHeader(5002)
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

func HandleCanvasOAuth2Response(w http.ResponseWriter, resp *http.Response, body string) {
	redirectTo, err := url.Parse(env.CanvasOAuth2SuccessURI)
	if err != nil {
		log.Fatal(err)
	}

	q := redirectTo.Query()
	q.Set("type", "canvas")

	if sc := resp.StatusCode; sc < 200 || sc > 399 {
		q.Set("error", "proxy_canvas_error")
		q.Set("error_source", "canvas_proxy")
		q.Set("canvas_status_code", fmt.Sprintf("%d", resp.StatusCode))

		// attempting JSON detection as Canvas sends text/html typed JSON with errors
		if string(body[:2]) == "{\"" ||
			strings.Contains(resp.Header.Get("content-type"), "application/json") {
			q.Set("body", body)
		} else {
			q.Set("body", "html_omitted")
		}
	} else {
		q.Set("canvas_response", body)
		q.Set("subdomain", env.CanvasOAuth2Subdomain)
	}

	redirectToURLString := fmt.Sprintf("%s?%s", env.CanvasOAuth2SuccessURI, q.Encode())

	SendRedirect(
		w,
		redirectToURLString,
	)
	return
}
