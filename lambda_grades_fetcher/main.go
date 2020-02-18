package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	s3Bucket            = "canvascbl-fetch-all-grades-logs"
	totalFetchErrorBody = "{\"error\":\"could not connect to host; no response received\"}"
)

var (
	sess = func() *session.Session {
		s, err := session.NewSession(&aws.Config{
			Region: aws.String(awsRegion),
		})
		if err != nil {
			panic(fmt.Errorf("error creating aws session: %w", err))
		}

		return s
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
	return time.Now().Format(time.RFC3339) + "-" + environment + ".json"
}

func HandleLambdaEvent() {
	ul := s3manager.NewUploader(sess)
	input := &s3manager.UploadInput{
		Bucket:      aws.String(s3Bucket),
		ContentType: aws.String("application/json"),
		Key:         aws.String(generateFilename()),
	}
	resp, err := client.Do(&req)
	if err != nil {
		if resp == nil {
			input.Body = strings.NewReader(totalFetchErrorBody)
		} else {
			input.Body = resp.Body
		}
	} else {
		input.Body = resp.Body
	}

	_, err = ul.Upload(input)
	if err != nil {
		// oh well, whatever.
		return
	}
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
