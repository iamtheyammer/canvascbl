package main

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/iamtheyammer/canvascbl/backend/src/admin"
	"github.com/iamtheyammer/canvascbl/backend/src/canvasapis"
	"github.com/iamtheyammer/canvascbl/backend/src/checkout"
	"github.com/iamtheyammer/canvascbl/backend/src/env"
	"github.com/iamtheyammer/canvascbl/backend/src/gradesapi"
	"github.com/iamtheyammer/canvascbl/backend/src/oauth2"
	"github.com/iamtheyammer/canvascbl/backend/src/plus"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"github.com/stripe/stripe-go"
	"log"
	"net/http"
	"time"
)

type MiddlewareRouter struct{}

var router = getRouter()

type notFound struct {
	Message []byte
}

func (nf notFound) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write(nf.Message)
	return
}

func getRouter() *httprouter.Router {
	// init http server

	router := &httprouter.Router{
		RedirectTrailingSlash:  true,
		HandleMethodNotAllowed: true,
		HandleOPTIONS:          true,
		NotFound: notFound{
			Message: []byte("404 not found - go.canvascbl.com/docs"),
		},
	}

	router.GET("/api/canvas/oauth2/request", canvasapis.OAuth2RequestHandler)
	router.GET("/api/canvas/oauth2/response", canvasapis.OAuth2ResponseHandler)

	router.GET("/api/checkout/session", checkout.CreateCheckoutSessionHandler)
	router.GET("/api/checkout/products", checkout.ListProductsHandler)
	router.POST("/api/checkout/redeem", checkout.RedeemHandler)
	router.GET("/api/checkout/subscriptions", checkout.ListSubscriptionsHandler)
	router.DELETE("/api/checkout/subscriptions", checkout.CancelSubscriptionHandler)
	// stripe webhook handler
	router.POST("/api/checkout/webhook", checkout.StripeWebhookHandler)

	router.GET("/api/plus/session", plus.GetSessionInformationHandler)
	router.DELETE("/api/plus/session", plus.ClearSessionHandler)
	router.GET("/api/plus/courses/:courseID/avg", plus.GetAverageGradeForCourseHandler)
	router.GET("/api/plus/outcomes/:outcomeID/avg", plus.GetAverageOutcomeScoreHandler)
	router.GET("/api/plus/grades/previous", plus.GetPreviousGradesHandler)

	// Fairly public
	router.GET("/api/oauth2/auth", oauth2.AuthHandler)
	// Session only, will carefully watch calls
	router.PUT("/api/oauth2/consent", oauth2.ConsentHandler)
	router.GET("/api/oauth2/consent", oauth2.ConsentInfoHandler)
	// Fairly public
	router.POST("/api/oauth2/token", oauth2.TokenHandler)
	router.DELETE("/api/oauth2/token", oauth2.DeleteTokenHandler)
	// Session only
	router.GET("/api/oauth2/tokens", oauth2.TokensHandler)

	// no google needed for now but we could bring it back later.
	//router.GET("/api/google/oauth2/request", googleapis.OAuth2RequestHandler)
	//router.GET("/api/google/oauth2/response", googleapis.OAuth2ResponseHandler)

	router.POST("/api/admin/gift_cards", admin.GenerateGiftCardsHandler)

	/*
		API
	*/

	// grades
	router.GET("/api/v1/grades", gradesapi.GradesHandler)
	router.GET("/api/v1/grades/fetch_all", gradesapi.GradesForAllHandler)

	// courses
	router.GET("/api/v1/courses/:courseID/assignments", gradesapi.AssignmentsHandler)
	router.GET("/api/v1/courses/:courseID/outcome_alignments", gradesapi.AlignmentsHandler)
	router.PUT("/api/v1/courses/:courseID/hide", gradesapi.HideCourseHandler)
	router.DELETE("/api/v1/courses/:courseID/hide", gradesapi.ShowCourseHandler)

	// outcomes
	router.GET("/api/v1/outcomes/:outcomeID", gradesapi.OutcomeHandler)

	// notifications
	router.GET("/api/v1/notifications/types", gradesapi.ListNotificationTypesHandler)
	router.GET("/api/v1/notifications/settings", gradesapi.ListNotificationSettingsHandler)
	router.PUT("/api/v1/notifications/types/:notificationTypeID", gradesapi.PutNotificationSettingsHandler)
	router.DELETE("/api/v1/notifications/types/:notificationTypeID", gradesapi.DeleteNotificationSettingHandler)
	return router
}

func (_ MiddlewareRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// apply CORS headers
	w.Header().Set("Access-Control-Allow-Origin", env.ProxyAllowedCORSOrigins)
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, "+
		"X-Canvas-Token, "+
		"X-Canvas-Subdomain, "+
		"X-Session-String")
	w.Header().Set("Access-Control-Expose-Headers", "X-Canvas-Url, X-Canvas-Status-Code")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if env.Env == env.EnvironmentStaging {
		if r.URL.Path != "/api/checkout/webhook" {
			if ok := util.CloudflareAccessVerifier.HandlerMiddleware(w, r); ok {
				return
			}
		}
	}

	router.ServeHTTP(w, r)
}

func main() {
	if env.Env != env.EnvironmentDevelopment {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:         env.SentryDSN,
			Environment: string(env.Env),
		})

		if err != nil {
			panic(errors.Wrap(err, "error initializing sentry"))
		} else {
			fmt.Println("INFO: Successfully initialized Sentry.")
		}
	}

	mw := MiddlewareRouter{}

	if env.ProxyAllowedCORSOrigins == "*" {
		fmt.Println("WARN: Your CANVAS_PROXY_ALLOWED_CORS_ORIGINS env var is currently set to \"*\", " +
			"which will allow any site to make requests to this server.")
	}

	stripe.Key = env.StripeAPIKey

	fmt.Println(fmt.Sprintf("Canvas proxy running on %s", env.HTTPPort))

	// Close db
	defer util.DB.Close()

	// Flush sentry queue
	defer sentry.Flush(5 * time.Second)

	log.Fatal(http.ListenAndServe(env.HTTPPort, mw))
}
