package main

import (
	"fmt"
	"github.com/iamtheyammer/canvascbl/backend/src/canvasapis"
	"github.com/iamtheyammer/canvascbl/backend/src/checkout"
	"github.com/iamtheyammer/canvascbl/backend/src/env"
	"github.com/iamtheyammer/canvascbl/backend/src/plus"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/julienschmidt/httprouter"
	"github.com/stripe/stripe-go"
	"log"
	"net/http"
)

type MiddlewareRouter struct{}

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
	router.GET("/api/canvas/users/profile/self/observees", canvasapis.GetOwnObserveesHandler)
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
	// saves to db
	router.GET("/api/canvas/courses/:courseID/outcome_results", canvasapis.GetOutcomeResultsByCourseHandler)
	// saves grades and specific outcome scores to DB
	router.GET("/api/canvas/courses/:courseID/outcome_rollups", canvasapis.GetOutcomeRollupsByCourseHandler)
	// doesn't save to db but a possibility
	router.GET("/api/canvas/courses/:courseID/outcome_alignments", canvasapis.GetOutcomeAlignmentsByCourseHandler)

	router.GET("/api/canvas/oauth2/request", canvasapis.OAuth2RequestHandler)
	router.GET("/api/canvas/oauth2/response", canvasapis.OAuth2ResponseHandler)
	router.GET("/api/canvas/oauth2/refresh_token", canvasapis.OAuth2RefreshTokenHandler)
	router.DELETE("/api/canvas/oauth2/token", canvasapis.DeleteOAuth2TokenHandler)

	router.POST("/api/canvas/tokens", canvasapis.InsertCanvasTokenHandler)
	router.GET("/api/canvas/tokens", canvasapis.GetCanvasTokensHandler)
	router.DELETE("/api/canvas/tokens", canvasapis.DeleteCanvasTokenHandler)

	router.GET("/api/checkout/session", checkout.CreateCheckoutSessionHandler)
	router.GET("/api/checkout/products", checkout.ListProductsHandler)
	router.GET("/api/checkout/subscriptions", checkout.ListSubscriptionsHandler)
	router.DELETE("/api/checkout/subscriptions", checkout.CancelSubscriptionHandler)
	// stripe webhook handler
	router.POST("/api/checkout/webhook", checkout.StripeWebhookHandler)

	router.GET("/api/plus/session", plus.GetSessionInformationHandler)
	router.GET("/api/plus/courses/:courseID/avg", plus.GetAverageGradeForCourseHandler)
	router.GET("/api/plus/outcomes/:outcomeID/avg", plus.GetAverageOutcomeScoreHandler)
	router.GET("/api/plus/grades/previous", plus.GetPreviousGradesHandler)

	// no google needed for now but we could bring it back later.
	//router.GET("/api/google/oauth2/request", googleapis.OAuth2RequestHandler)
	//router.GET("/api/google/oauth2/response", googleapis.OAuth2ResponseHandler)

	return router
}

func (_ MiddlewareRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// apply CORS headers
	w.Header().Set("Access-Control-Allow-Origin", env.ProxyAllowedCORSOrigins)
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, "+
		"X-Canvas-Token, "+
		"X-Canvas-Subdomain, "+
		"X-Session-String")
	w.Header().Set("Access-Control-Expose-Headers", "X-Canvas-Url, X-Canvas-Status-Code")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if env.Env == env.EnvironmentStaging {
		if ok := util.CloudflareAccessVerifier.HandlerMiddleware(w, r); ok {
			return
		}
	}

	router.ServeHTTP(w, r)
}

func main() {
	mw := MiddlewareRouter{}

	if env.ProxyAllowedSubdomains[0] == "*" {
		fmt.Println("WARN: Your CANVAS_PROXY_ALLOW_SUBDOMAINS env var is currently set to \"*\", " +
			"which will allow anyone to use this server as a proxy server.")
	}

	if env.ProxyAllowedCORSOrigins == "*" {
		fmt.Println("WARN: Your CANVAS_PROXY_ALLOWED_CORS_ORIGINS env var is currently set to \"*\", " +
			"which will allow any site to make requests to this server.")
	}

	stripe.Key = env.StripeAPIKey

	fmt.Println(fmt.Sprintf("Canvas proxy running on %s", env.HTTPPort))

	// Close db
	defer util.DB.Close()

	log.Fatal(http.ListenAndServe(env.HTTPPort, mw))
}
