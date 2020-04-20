package gradesapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/iamtheyammer/canvascbl/backend/src/env"
	"github.com/tomnomnom/linkheader"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const canvasPerPage = "100"

var (
	canvasErrorNoErrors                             = errors.New("a non-200 status code was received but no errors were present")
	canvasErrorUnknownError                         = errors.New("a non-200 status code was received from Canvas, but the error is unknown")
	canvasErrorInvalidAccessTokenError              = errors.New("the Canvas access token is invalid")
	canvasErrorInsufficientScopesOnAccessTokenError = errors.New("there are insufficient scopes on the Canvas access token")

	canvasOAuth2ErrorRefreshTokenNotFound = errors.New("the specified refresh token was not found")

	handleRequestWithTokenRefreshMutex = sync.RWMutex{}
	lockedTokens                       = map[uint64]struct{}{}
)

var proxyURL, _ = url.Parse("http://localhost:8888")
var httpClient = http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}

//var httpClient = http.Client{}

type requestDetails struct {
	// TokenID is the database ID of the token
	TokenID uint64
	// Token represents the user's Canvas token
	Token string
	// RefreshToken represents the user's Canvas refresh token
	RefreshToken string
}

/*
handleRequestWithTokenRefresh takes in a task function along with your requestDetails,
runs the task, and, if necessary, refreshes the token and retries the task.

It is also safe for concurrent use-- meaning this is the only way to make canvas requests.

Your task function should take requestDetails and return an error. This error, if wrapped,
should be from fmt.Errorf using the %w verb. For things like requestDetails and parameters,
they should be scoped in from your function.

*/
func handleRequestWithTokenRefresh(task func(rd *requestDetails) error, rd *requestDetails, userID uint64) (requestDetails, error) {
	err := task(rd)
	if err != nil {
		if errors.Is(err, canvasErrorInvalidAccessTokenError) {
			// we need to use the refresh token

			shouldRefresh := false

			handleRequestWithTokenRefreshMutex.Lock()
			// if the token isn't being refreshed
			if _, ok := lockedTokens[rd.TokenID]; !ok {
				// mark that it is
				lockedTokens[rd.TokenID] = struct{}{}
				shouldRefresh = true
			}
			handleRequestWithTokenRefreshMutex.Unlock()

			if shouldRefresh {
				// refresh the token
				refreshErr := rd.refreshAccessToken()
				if refreshErr != nil {
					return requestDetails{}, fmt.Errorf("error refreshing token id %d: %w", rd.TokenID, refreshErr)
				}
				// mark that it's no longer being refreshed
				handleRequestWithTokenRefreshMutex.Lock()
				delete(lockedTokens, rd.TokenID)
				handleRequestWithTokenRefreshMutex.Unlock()
			} else {
				// poll the map every 2ms for updates
				for {
					// check map
					handleRequestWithTokenRefreshMutex.RLock()
					if _, ok := lockedTokens[rd.TokenID]; ok {
						// still working, check back in 2ms
						handleRequestWithTokenRefreshMutex.RUnlock()
						time.Sleep(2 * time.Millisecond)
					} else {
						handleRequestWithTokenRefreshMutex.RUnlock()

						// get new token from db
						newRd, err := rdFromUserID(userID)
						if err != nil {
							return requestDetails{}, fmt.Errorf("error getting rd from canvas user id: %w", err)
						}

						rd = &newRd
						break
					}
				}
			}

			retryErr := task(rd)
			if retryErr != nil {
				return requestDetails{}, fmt.Errorf("error retrying task with refreshed token id %d: %w", rd.TokenID, retryErr)
			}
		} else {
			return requestDetails{}, fmt.Errorf("error in task from handleRequestWithTokenRefresh: %w", err)
		}
	}

	return *rd, nil
}

func getCanvasProfile(rd requestDetails, userID string) (*canvasUserProfileResponse, error) {
	var profile canvasUserProfileResponse
	_, err := makeCanvasGetRequest("api/v1/users/"+userID+"/profile", rd, &profile)
	if err != nil {
		return nil, fmt.Errorf("error getting canvas user profile: %w", err)
	}

	return &profile, nil
}

func getCanvasCourses(rd requestDetails) (*canvasCoursesResponse, error) {
	var courses canvasCoursesResponse
	_, err := makeCanvasGetRequest("api/v1/courses?per_page=100&enrollment_state=active&include[]=total_scores&include[]=observed_users", rd, &courses)
	if err != nil {
		return nil, fmt.Errorf("error getting canvas courses: %w", err)
	}

	return &courses, nil
}

func getCanvasUserObservees(rd requestDetails, userID string) (*canvasUserObserveesResponse, error) {
	var observees canvasUserObserveesResponse
	_, err := makeCanvasGetRequest("api/v1/users/"+userID+"/observees", rd, &observees)
	if err != nil {
		return nil, fmt.Errorf("error getting canvas user observees: %w", err)
	}

	return &observees, nil
}

// getCanvasOutcomeAlignments is not currently needed, but may be useful in the future.
//func getCanvasOutcomeAlignments(rd requestDetails, courseID string, studentID string) (*canvasOutcomeAlignmentsResponse, error) {
//	var alignments canvasOutcomeAlignmentsResponse
//	_, err := makeCanvasGetRequest(
//		"api/v1/courses/"+courseID+"/outcome_alignments?student_id="+studentID,
//		rd,
//		&alignments,
//	)
//	if err != nil {
//		return nil, fmt.Errorf("error getting canvas outcome alignments for course %s for student %s: %w", courseID, studentID, err)
//	}
//
//	return &alignments, nil
//}

func proxyCanvasOutcomeAlignments(rd requestDetails, courseID string, studentID string) (*http.Response, error) {
	resp, err := proxyCanvasGetRequest("api/v1/courses/"+courseID+"/outcome_alignments?per_page=100&student_id="+studentID, rd)
	if err != nil {
		return nil, fmt.Errorf("error getting canvas outcome alignments for course %s: %w", courseID, err)
	}

	return resp, nil
}

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
//		"api/v1/courses/"+courseID+"/outcome_rollups"+q.Encode(),
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

// getCanvasOutcomeResults gets outcome results. Paginated.
func getCanvasOutcomeResults(rd requestDetails, courseID string, userIDs []string) (*canvasOutcomeResultsResponse, error) {
	q := url.Values{}
	q.Add("per_page", canvasPerPage)
	for _, id := range userIDs {
		q.Add("user_ids[]", id)
	}

	var (
		allResults canvasOutcomeResultsResponse
		u          = "api/v1/courses/" + courseID + "/outcome_results?" + q.Encode()
	)

	for {
		var results canvasOutcomeResultsResponse
		resp, err := makeCanvasGetRequest(
			u,
			rd,
			&results,
		)
		if err != nil {
			return nil, fmt.Errorf("error getting canvas outcome results for course "+
				"%s and for students %v: %w", courseID, userIDs, err)
		}

		allResults.OutcomeResults = append(allResults.OutcomeResults, results.OutcomeResults...)

		nu := nextPageUrl(resp.Header.Get("link"))
		if nu == nil {
			return &allResults, nil
		} else {
			u = *nu
		}
	}

}

func getCanvasCourseAssignments(rd requestDetails, courseID string, assignmentIDs []string) (*canvasAssignmentsResponse, error) {
	q := url.Values{}
	q.Add("per_page", canvasPerPage)
	for _, aID := range assignmentIDs {
		q.Add("assignment_ids[]", aID)
	}

	var (
		allAssignments canvasAssignmentsResponse
		u              = "api/v1/courses/" + courseID + "/assignments?" + q.Encode()
	)

	for {
		var assignments canvasAssignmentsResponse
		resp, err := makeCanvasGetRequest(u, rd, &assignments)
		if err != nil {
			return nil, fmt.Errorf("error getting canvas assignments for course %s: %w", courseID, err)
		}

		allAssignments = append(allAssignments, assignments...)

		nu := nextPageUrl(resp.Header.Get("link"))
		if nu == nil {
			return &allAssignments, nil
		} else {
			u = *nu
		}
	}
}

// getCanvasCourseEnrollments gets all enrollments for a course. This request can be performed by
// any user, but only teachers get grades.
func getCanvasCourseEnrollments(rd requestDetails, courseID string) (*canvasEnrollmentsResponse, error) {
	q := url.Values{}
	q.Add("per_page", canvasPerPage)
	q.Add("type[]", "StudentEnrollment")
	q.Add("state[]", "active")
	q.Add("include[]", "avatar_url")

	var (
		allEnrollments canvasEnrollmentsResponse
		u              = "api/v1/courses/" + courseID + "/enrollments?" + q.Encode()
	)

	for {
		var enrollments canvasEnrollmentsResponse
		resp, err := makeCanvasGetRequest(u, rd, &enrollments)
		if err != nil {
			return nil, fmt.Errorf("error getting canvas enrollments for course %s: %w", courseID, err)
		}

		allEnrollments = append(allEnrollments, enrollments...)

		nu := nextPageUrl(resp.Header.Get("link"))
		if nu == nil {
			return &allEnrollments, nil
		} else {
			u = *nu
		}
	}
}

func getCanvasOutcome(rd requestDetails, outcomeID string) (*canvasOutcomeResponse, error) {
	var outcome canvasOutcomeResponse
	_, err := makeCanvasGetRequest("api/v1/outcomes/"+outcomeID, rd, &outcome)
	if err != nil {
		return nil, fmt.Errorf("error getting canvas outcome %s: %w", outcomeID, err)
	}

	return &outcome, err
}

func getTokenFromAuthorizationCode(code string) (*canvasTokenGrantResponse, error) {
	q := url.Values{}
	q.Set("grant_type", "authorization_code")
	q.Set("client_id", env.CanvasOAuth2ClientID)
	q.Set("client_secret", env.CanvasOAuth2ClientSecret)
	q.Set("code", code)

	var grantResp canvasTokenGrantResponse
	_, err := makeCanvasRequest(
		"login/oauth2/token?"+q.Encode(),
		http.MethodPost,
		nil,
		requestDetails{},
		&grantResp,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting a canvas access token from a refresh token: %w", err)
	}

	return &grantResp, nil
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

func categorizeCanvasError(err canvasErrorArrayResponse, resp *http.Response) error {
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

func categorizeCanvasOAuth2Error(err canvasOAuth2ErrorResponse, resp *http.Response) error {
	switch strings.ToLower(err.ErrorDescription) {
	case "refresh_token not found":
		return canvasOAuth2ErrorRefreshTokenNotFound
	default:
		return fmt.Errorf("oauth2 error from Canvas with status code %d: %w", resp.StatusCode, canvasErrorUnknownError)
	}
}

// Convenience method for makeCanvasRequest with no body and the method set to get
func makeCanvasGetRequest(path string, rd requestDetails, bodyDestination interface{}) (*http.Response, error) {
	return makeCanvasRequest(path, http.MethodGet, nil, rd, bodyDestination)
}

// Convenience method for makeCanvasRequest with the method set to post
func makeCanvasPostRequest(path string, body io.Reader, rd requestDetails, bodyDestination interface{}) (*http.Response, error) {
	return makeCanvasRequest(path, http.MethodPost, body, rd, bodyDestination)
}

func nextPageUrl(link string) *string {
	if len(link) < 1 {
		return nil
	}

	links := linkheader.Parse(link)
	for _, l := range links {
		if l.Rel == "next" {
			return &l.URL
		}
	}

	return nil
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
	// this system allows for the Link header so full URLs can be passed in
	fURL := "https://" + env.CanvasDomain + "/"
	if strings.HasPrefix(path, fURL) {
		fURL = path
	} else {
		fURL += path
	}

	req, err := http.NewRequest(method, fURL, body)
	if err != nil {
		return nil, fmt.Errorf("error creating http request: %w", err)
	}

	if len(rd.Token) > 0 {
		req.Header.Add("Authorization", "Bearer "+rd.Token)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making an http request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if strings.HasPrefix(path, "login/oauth2") {
			var oauth2Err canvasOAuth2ErrorResponse
			err = json.NewDecoder(resp.Body).Decode(&oauth2Err)
			if err != nil {
				return nil, fmt.Errorf(
					"error decoding into canvasOAuth2Error (canvas status code %d): %w",
					resp.StatusCode,
					err)
			}

			return nil, fmt.Errorf(
				"oauth2 error from canvas (canvas status code %d): %w",
				resp.StatusCode,
				categorizeCanvasOAuth2Error(oauth2Err, resp),
			)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf(
				"error reading the canvas error body (canvas status code %d): %w",
				resp.StatusCode,
				err)
		}

		var canvasArrayErr canvasErrorArrayResponse
		err = json.Unmarshal(body, &canvasArrayErr)
		if err != nil {
			return nil, fmt.Errorf(
				"error decoding into canvasArrayErr (canvas status code %d): %w",
				resp.StatusCode,
				err)
		}

		if len(canvasArrayErr.Errors) < 1 {
			var canvasErr canvasErrorResponse
			err = json.Unmarshal(body, &canvasErr)
			if err != nil {
				return nil, fmt.Errorf(
					"error decoding into canvasErr (canvas status code %d): %w",
					resp.StatusCode,
					err)
			}

			if len(canvasErr.Error) > 0 {
				canvasArrayErr = canvasErr.toCanvasErrorArrayResponse()
			}
		}

		return nil, fmt.Errorf(
			"error from canvas (canvas status code %d): %w",
			resp.StatusCode,
			categorizeCanvasError(canvasArrayErr, resp),
		)
	}

	err = json.NewDecoder(resp.Body).Decode(&bodyDestination)
	if err != nil {
		return nil, fmt.Errorf("error decoding into bodyDestination: %w", err)
	}

	return resp, nil
}

// proxyCanvasGetRequest expects you to read resp.Body. So, it doesn't close the body.
// REMEMBER TO CLOSE IT!
func proxyCanvasGetRequest(path string, rd requestDetails) (*http.Response, error) {
	// this system allows for the Link header so full URLs can be passed in
	fURL := "https://" + env.CanvasDomain + "/"
	if !strings.HasPrefix(path, fURL) {
		fURL += path
	}
	req, err := http.NewRequest(http.MethodGet, fURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating http request: %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+rd.Token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making an http request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf(
				"error reading the canvas error body (canvas status code %d): %w",
				resp.StatusCode,
				err)
		}

		var canvasArrayErr canvasErrorArrayResponse
		err = json.Unmarshal(body, &canvasArrayErr)
		if err != nil {
			return nil, fmt.Errorf(
				"error decoding into canvasArrayErr (canvas status code %d): %w",
				resp.StatusCode,
				err)
		}

		if len(canvasArrayErr.Errors) < 1 {
			var canvasErr canvasErrorResponse
			err = json.Unmarshal(body, &canvasErr)
			if err != nil {
				return nil, fmt.Errorf(
					"error decoding into canvasErr (canvas status code %d): %w",
					resp.StatusCode,
					err)
			}

			if len(canvasErr.Error) > 0 {
				canvasArrayErr = canvasErr.toCanvasErrorArrayResponse()
			}
		}

		return nil, fmt.Errorf(
			"error from canvas (canvas status code %d): %w",
			resp.StatusCode,
			categorizeCanvasError(canvasArrayErr, resp),
		)
	}

	return resp, nil
}
