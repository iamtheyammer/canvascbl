package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"net/http"
	"net/url"
	"time"
)

const (
	s3Bucket = "canvascbl-fetch-all-grades-logs"
)

var (
	ul = func() *s3manager.Uploader {
		s, err := session.NewSession(&aws.Config{
			Region: aws.String(awsRegion),
		})
		if err != nil {
			panic(fmt.Errorf("error creating aws session: %w", err))
		}

		upl := s3manager.NewUploader(s)
		return upl
	}()
	client = http.Client{}
	apiURL = func() *url.URL {
		parsedURL, err := url.Parse(apiURLEnv)
		if err != nil {
			panic(fmt.Errorf("error parsing API_URL: %w", err))
		}

		return parsedURL
	}()
	req = http.Request{
		Method: "GET",
		URL:    apiURL,
		Header: http.Header{
			"X-CanvasCBL-Script-Key": []string{scriptKey},
		},
	}
)

func generateFilename() string {
	return "lgf-error-" + time.Now().Format(time.RFC3339) + "-" + environment + ".json"
}

func HandleLambdaEvent() error {
	// create a new context
	ctx := context.Background()
	// add a 1 second deadline (enough to connect and process the request)
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*1))
	// add context to request
	r := req.WithContext(ctx)
	// make the request
	resp, err := client.Do(r)
	// if an error occurs and resp != nil upload error to s3
	if err != nil && resp != nil {
		input := &s3manager.UploadInput{
			Bucket:      aws.String(s3Bucket),
			ContentType: aws.String("application/json"),
			Key:         aws.String(generateFilename()),
			Body:        resp.Body,
		}

		_, err := ul.Upload(input)
		if err != nil {
			return fmt.Errorf("error uploading to s3: %w", err)
		}
	}

	// you're supposed to always call cancel
	cancel()

	return nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
