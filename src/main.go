package main

import (
	"fmt"
	"github.com/iamtheyammer/canvascbl/backend/src/canvasapis"
	"github.com/iamtheyammer/canvascbl/backend/src/env"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

var router = getRouter()

func getRouter() *httprouter.Router {
	// init http server

	router := &httprouter.Router{
		RedirectTrailingSlash:  true,
		HandleMethodNotAllowed: true,
		HandleOPTIONS:          true,
	}

	// saves to db
	router.GET("/api/canvas/outcomes/:outcomeID", canvasapis.GetOutcomeByIDHandler)
	// saves to db
	router.GET("/api/canvas/users/profile/self", canvasapis.GetOwnUserProfileHandler)
	// saves to db
	router.GET("/api/canvas/courses", canvasapis.GetCoursesHandler)
	// saves to db
	router.GET("/api/canvas/courses/:courseID/assignments", canvasapis.GetAssignmentsByCourseHandler)
	// no db needed
	router.GET("/api/canvas/courses/:courseID/outcome_groups", canvasapis.GetOutcomesByCourseHandler)
	router.GET(
		"/api/canvas/courses/:courseID/outcome_groups/:outcomeGroupID/outcomes",
		canvasapis.GetOutcomesByCourseAndOutcomeGroupHandler,
	)
	// TODO: save to db
	router.GET("/api/canvas/courses/:courseID/outcome_results", canvasapis.GetOutcomeResultsByCourseHandler)
	// saves grades to DB, TODO: save specific outcome scores
	router.GET("/api/canvas/courses/:courseID/outcome_rollups", canvasapis.GetOutcomeRollupsByCourseHandler)

	router.GET("/api/canvas/oauth2/request", canvasapis.OAuth2RequestHandler)
	router.GET("/api/canvas/oauth2/response", canvasapis.OAuth2ResponseHandler)
	router.GET("/api/canvas/oauth2/refresh_token", canvasapis.OAuth2RefreshTokenHandler)
	router.DELETE("/api/canvas/oauth2/token", canvasapis.DeleteOAuth2TokenHandler)

	return router
}

type MiddlewareRouter map[string]string

func (_ MiddlewareRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// apply CORS headers
	w.Header().Set("Access-Control-Allow-Origin", env.ProxyAllowedCORSOrigins)
	w.Header().Set("Access-Control-Allow-Methods", "GET, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "X-Canvas-Token, X-Canvas-Subdomain")
	w.Header().Set("Access-Control-Expose-Headers", "X-Canvas-Url, X-Canvas-Status-Code")

	router.ServeHTTP(w, r)
}

func main() {
	mw := make(MiddlewareRouter)

	if env.ProxyAllowedSubdomains[0] == "*" {
		fmt.Println("WARN: Your CANVAS_PROXY_ALLOW_SUBDOMAINS env var is currently set to \"*\", " +
			"which will allow anyone to use this server as a proxy server.")
	}

	if env.ProxyAllowedCORSOrigins == "*" {
		fmt.Println("WARN: Your CANVAS_PROXY_ALLOWED_CORS_ORIGINS env var is currently set to \"*\", " +
			"which will allow any site to make requests to this server.")
	}

	fmt.Println(fmt.Sprintf("Canvas proxy running on %s", env.HTTPPort))

	// Close db
	defer util.DB.Close()

	log.Fatal(http.ListenAndServe(env.HTTPPort, mw))
}
