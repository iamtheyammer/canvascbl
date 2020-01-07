package courses

import (
	"fmt"
	"github.com/iamtheyammer/canvascbl/backend/src/req"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"net/http"
	"net/url"
)

// Get gets all of a user's courses
func Get(rd *util.RequestDetails) (*http.Response, string, error) {
	u := fmt.Sprintf(
		"https://%s.instructure.com/api/v1/courses?per_page=100",
		rd.Subdomain,
	)
	return req.MakeAuthenticatedGetRequest(u, rd.Token)
}

// GetOutcomesByCourse gets all outcomes for a specific course
func GetOutcomesByCourse(rd *util.RequestDetails, courseID string) (*http.Response, string, error) {
	u := fmt.Sprintf(
		"https://%s.instructure.com/api/v1/courses/%s/outcome_groups",
		rd.Subdomain,
		courseID,
	)
	return req.MakeAuthenticatedGetRequest(u, rd.Token)
}

// GetOutcomesByCourseAndOutcomeGroup gets all outcomes in a course's outcome group
func GetOutcomesByCourseAndOutcomeGroup(
	rd *util.RequestDetails,
	courseID string,
	outcomeGroupID string,
) (*http.Response, string, error) {
	u := fmt.Sprintf(
		"https://%s.instructure.com/api/v1/courses/%s/outcome_groups/%s/outcomes?per_page=100",
		rd.Subdomain,
		courseID,
		outcomeGroupID,
	)
	return req.MakeAuthenticatedGetRequest(u, rd.Token)
}

// GetOutcomeResultsByCourse gets outcome results for the specified course
func GetOutcomeResultsByCourse(
	rd *util.RequestDetails,
	courseID string,
	userIDs []string,
	include string,
) (*http.Response, string, error) {
	u := url.URL{
		Scheme: "https",
		Host:   fmt.Sprintf("%s.instructure.com", rd.Subdomain),
		Path:   fmt.Sprintf("/api/v1/courses/%s/outcome_results", courseID),
	}

	q := u.Query()
	for _, v := range userIDs {
		q.Add("user_ids[]", v)
	}

	if len(include) > 0 {
		q.Add("include[]", include)
	}

	// max = 100
	q.Add("per_page", "100")
	u.RawQuery = q.Encode()

	return req.MakeAuthenticatedGetRequest(u.String(), rd.Token)
}

// GetOutcomeRollupsByCourse gets outcome rollups for a specific course
func GetOutcomeRollupsByCourse(
	rd *util.RequestDetails,
	courseID string,
	userIDs []string,
	include string,
) (*http.Response, string, error) {
	u := url.URL{
		Scheme: "https",
		Host:   fmt.Sprintf("%s.instructure.com", rd.Subdomain),
		Path:   fmt.Sprintf("/api/v1/courses/%s/outcome_rollups", courseID),
	}

	q := u.Query()
	for _, v := range userIDs {
		q.Add("user_ids[]", v)
	}

	if len(include) > 0 {
		q.Add("include[]", include)
	}

	// max = 100
	q.Add("per_page", "100")
	u.RawQuery = q.Encode()

	return req.MakeAuthenticatedGetRequest(u.String(), rd.Token)
}

// GetAssignmentsByCourse gets all assignments for a specified course
func GetAssignmentsByCourse(
	rd *util.RequestDetails,
	courseID string,
	include string,
) (*http.Response, string, error) {
	u := fmt.Sprintf(
		"https://%s.instructure.com/api/v1/courses/%s/assignments?per_page=100",
		rd.Subdomain,
		courseID,
	)

	if len(include) > 1 {
		u = fmt.Sprintf("%s&include[]=%s", u, include)
	}

	return req.MakeAuthenticatedGetRequest(u, rd.Token)
}

func GetOutcomeAlignmentsByCourse(
	rd *util.RequestDetails,
	courseID string,
	userID string,
) (*http.Response, string, error) {
	u := fmt.Sprintf(
		"https://%s.instructure.com/api/v1/courses/%s/outcome_alignments?student_id=%s",
		rd.Subdomain,
		courseID,
		userID,
	)

	return req.MakeAuthenticatedGetRequest(u, rd.Token)
}
