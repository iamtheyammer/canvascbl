package gradesapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/iamtheyammer/canvascbl/backend/src/env"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const canvasPerPage = "100"

var (
	canvasErrorNoErrors                             = errors.New("a non-200 status code was received but no errors were present")
	canvasErrorUnknownError                         = errors.New("a non-200 status code was received from Canvas, but the error is unknown")
	canvasErrorInvalidAccessTokenError              = errors.New("the Canvas access token is invalid")
	canvasErrorInsufficientScopesOnAccessTokenError = errors.New("there are insufficient scopes on the Canvas access token")
)

var proxyURL, _ = url.Parse("http://localhost:8888")
var tr = &http.Transport{
	Proxy: http.ProxyURL(proxyURL),
}

var httpClient = http.Client{Transport: tr}

type requestDetails struct {
	// TokenID is the database ID of the token
	TokenID uint64
	// Token represents the user's Canvas token
	Token string
	// RefreshToken represents the user's Canvas refresh token
	RefreshToken string
	// Subdomain represents the user's Canvas Subdomain
	Subdomain string
}

func getCanvasProfile(rd requestDetails, userID string) (*canvasUserProfileResponse, error) {
	var profile canvasUserProfileResponse
	_, err := makeCanvasGetRequest("users/"+userID+"/profile", rd, &profile)
	if err != nil {
		return nil, fmt.Errorf("error getting canvas user profile: %w", err)
	}

	return &profile, nil
}

func getCanvasCourses(rd requestDetails) (*canvasCoursesResponse, error) {
	var courses canvasCoursesResponse
	_, err := makeCanvasGetRequest("courses", rd, &courses)
	if err != nil {
		return nil, fmt.Errorf("error getting canvas courses: %w", err)
	}

	return &courses, nil
}

func getCanvasUserObservees(rd requestDetails, userID string) (*canvasUserObserveesResponse, error) {
	var observees canvasUserObserveesResponse
	_, err := makeCanvasGetRequest("users/"+userID+"/observees", rd, &observees)
	if err != nil {
		return nil, fmt.Errorf("error getting canvas user observees: %w", err)
	}

	return &observees, nil
}

// getCanvasOutcomeAlignments is not currently needed, but may be useful in the future.
//func getCanvasOutcomeAlignments(rd requestDetails, courseID string, studentID string) (*canvasOutcomeAlignmentsResponse, error) {
//	var alignments canvasOutcomeAlignmentsResponse
//	_, err := makeCanvasGetRequest(
//		"courses/"+courseID+"/outcome_alignments?student_id="+studentID,
//		rd,
//		&alignments,
//	)
//	if err != nil {
//		return nil, fmt.Errorf("error getting canvas outcome alignments for course %s for student %s: %w", courseID, studentID, err)
//	}
//
//	return &alignments, nil
//}

//func proxyCanvasOutcomeAlignments(rd requestDetails, courseID string, studentID string) (*http.Response, error) {
//	resp, err := proxyCanvasGetRequest("courses/"+courseID+"/outcome_alignments?student_id="+studentID, rd)
//	if err != nil {
//		return nil, fmt.Errorf("error getting canvas outcome alignments for course %s: %w", courseID, err)
//	}
//
//	return resp, nil
//}

// getCanvasOutcomeRollups is currently deprecated but still here because it may be useful in the future.
//func getCanvasOutcomeRollups(rd requestDetails, courseID string, userIDs []string) (*canvasOutcomeRollupsResponse, error) {
//	q := url.Values{}
//	q.Add("per_page", canvasPerPage)
//	for _, id := range userIDs {
//		q.Add("user_ids[]", id)
//	}
//
//	var rollups canvasOutcomeRollupsResponse
//	_, err := makeCanvasGetRequest(
//		"courses/"+courseID+"/outcome_rollups"+q.Encode(),
//		rd,
//		&rollups,
//	)
//	if err != nil {
//		return nil, fmt.Errorf("error getting canvas outcome rollups for course "+
//			"%s and for students %v: %w", courseID, userIDs, err)
//	}
//
//	return &rollups, nil
//}

func getCanvasOutcomeResults(rd requestDetails, courseID string, userIDs []string) (*canvasOutcomeResultsResponse, error) {
	q := url.Values{}
	q.Add("per_page", canvasPerPage)
	for _, id := range userIDs {
		q.Add("user_ids[]", id)
	}

	var results canvasOutcomeResultsResponse
	_, err := makeCanvasGetRequest(
		"courses/"+courseID+"/outcome_results?"+q.Encode(),
		rd,
		&results,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting canvas outcome results for course "+
			"%s and for students %v: %w", courseID, userIDs, err)
	}

	return &results, nil
}

func getCanvasCourseAssignments(rd requestDetails, courseID string) (*canvasAssignmentsResponse, error) {
	q := url.Values{}
	q.Add("per_page", canvasPerPage)

	var assignments canvasAssignmentsResponse
	_, err := makeCanvasGetRequest("courses/"+courseID+"/assignments?"+q.Encode(), rd, &assignments)
	if err != nil {
		return nil, fmt.Errorf("error getting canvas assignments for course %s: %w", courseID, err)
	}

	return &assignments, nil
}

func getCanvasOutcome(rd requestDetails, outcomeID string) (*canvasOutcomeResponse, error) {
	var outcome canvasOutcomeResponse
	_, err := makeCanvasGetRequest("outcomes/"+outcomeID, rd, &outcome)
	if err != nil {
		return nil, fmt.Errorf("error getting canvas outcome %s: %w", outcomeID, err)
	}

	return &outcome, err
}

func getTokenFromRefreshToken(rd requestDetails) (*canvasRefreshTokenResponse, error) {
	q := url.Values{}
	q.Set("client_id", env.CanvasOAuth2ClientID)
	q.Set("client_secret", env.CanvasOAuth2ClientSecret)
	q.Set("grant_type", "refresh_token")
	q.Set("refresh_token", rd.RefreshToken)

	var rtResponse canvasRefreshTokenResponse
	_, err := makeCanvasRequest(
		"login/oauth2/token?"+q.Encode(),
		http.MethodPost,
		nil,
		rd,
		&rtResponse,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting a canvas access token from a refresh token: %w", err)
	}

	return &rtResponse, nil
}

func categorizeCanvasError(err canvasErrorResponse, resp *http.Response) error {
	if len(err.Errors) < 1 {
		return canvasErrorNoErrors
	}

	switch strings.ToLower(err.Errors[0].Message) {
	case "invalid access token.":
		return canvasErrorInvalidAccessTokenError
	case "insufficient scopes on access token.":
		return canvasErrorInsufficientScopesOnAccessTokenError
	default:
		return fmt.Errorf("error from Canvas with status code %d: %w", resp.StatusCode, canvasErrorUnknownError)
	}
}

// Convenience method for makeCanvasRequest with no body and the method set to get
func makeCanvasGetRequest(path string, rd requestDetails, bodyDestination interface{}) (*http.Response, error) {
	return makeCanvasRequest("api/v1/"+path, http.MethodGet, nil, rd, bodyDestination)
}

// Convenience method for makeCanvasRequest with the method set to post
func makeCanvasPostRequest(path string, body io.Reader, rd requestDetails, bodyDestination interface{}) (*http.Response, error) {
	return makeCanvasRequest("api/v1/"+path, http.MethodPost, body, rd, bodyDestination)
}

// makeCanvasGetRequest will WRITE TO YOUR bodyDestination.
// While it returns an http.Response, it is foolish to read the body because
// it's already in your bodyDestination. Ensure that bodyDestination is
// a POINTER to a struct with json labels.
func makeCanvasRequest(
	path string,
	method string,
	body io.Reader,
	rd requestDetails,
	bodyDestination interface{},
) (*http.Response, error) {
	fURL := "https://" + rd.Subdomain + ".instructure.com/" + path
	req, err := http.NewRequest(method, fURL, body)
	if err != nil {
		return nil, fmt.Errorf("error creating http request: %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+rd.Token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making an http request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var canvasErr canvasErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&canvasErr)
		if err != nil {
			return nil, fmt.Errorf(
				"error decoding into canvasErr (canvas status code %d): %w",
				resp.StatusCode,
				err)
		}

		return nil, fmt.Errorf(
			"error from canvas (canvas status code %d): %w",
			resp.StatusCode,
			categorizeCanvasError(canvasErr, resp),
		)
	}

	err = json.NewDecoder(resp.Body).Decode(&bodyDestination)
	if err != nil {
		return nil, fmt.Errorf("error decoding into bodyDestination: %w", err)
	}

	return resp, nil
}

//func proxyCanvasGetRequest(path string, rd requestDetails) (*http.Response, error) {
//	fURL := "https://" + rd.Subdomain + ".instructure.com/" + path
//	req, err := http.NewRequest(http.MethodGet, fURL, nil)
//	if err != nil {
//		return nil, fmt.Errorf("error creating http request: %w", err)
//	}
//
//	req.Header.Add("Authorization", "Bearer "+rd.Token)
//
//	resp, err := httpClient.Do(req)
//	if err != nil {
//		return nil, fmt.Errorf("error making an http request: %w", err)
//	}
//
//	defer resp.Body.Close()
//
//	if resp.StatusCode != http.StatusOK {
//		var canvasErr canvasErrorResponse
//		err = json.NewDecoder(resp.Body).Decode(&canvasErr)
//		if err != nil {
//			return nil, fmt.Errorf(
//				"error decoding into canvasErr (canvas status code %d): %w",
//				resp.StatusCode,
//				err)
//		}
//
//		return nil, fmt.Errorf(
//			"error from canvas (canvas status code %d): %w",
//			resp.StatusCode,
//			categorizeCanvasError(canvasErr, resp),
//		)
//	}
//
//	return resp, nil
//}
