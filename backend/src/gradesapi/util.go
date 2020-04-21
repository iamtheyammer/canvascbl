package gradesapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/canvas_tokens"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const spring20DLEnrollmentTermID = 3

// refreshAccessToken attempts to use the refresh token to get a new access token.
// You should check the returned error and handle it accordingly.
func (rd *requestDetails) refreshAccessToken() error {
	newToken, err := getTokenFromRefreshToken(*rd)
	if err != nil {
		if errors.Is(err, canvasErrorInvalidAccessTokenError) {
			// there is something wrong with the token
			deleteErr := canvas_tokens.Delete(db, &canvas_tokens.DeleteRequest{RefreshToken: rd.RefreshToken})
			if deleteErr != nil {
				return fmt.Errorf("error deleting a canvas token: %w", deleteErr)
			}

			return fmt.Errorf("error getting a new access token from a refresh token: %w", err)
		}

		return fmt.Errorf("error getting a newProfile: %w", err)
	}

	newTokenExp := time.Now().UTC().Add(time.Duration(newToken.ExpiresIn) * time.Second)
	err = canvas_tokens.UpdateFromRefreshToken(db, rd.TokenID, newToken.AccessToken, &newTokenExp)
	if err != nil {
		return fmt.Errorf("error updating a canvas token: %w", err)
	}

	rd.Token = newToken.AccessToken
	return nil
}

// getGradedUsersAndValidCourses gets graded users (users enrolled as a student in a course)
// and a list of courses that have either an observer enrollment or a student enrollment, as any other
// enrollment type can't get grades. Valid courses also have not ended.
func getGradedUsersAndValidCourses(courses *[]canvasCourse) (map[uint64][]uint64, *[]canvasCourse) {
	var (
		gradedUsers  = map[uint64][]uint64{}
		validCourses []canvasCourse
	)

	for _, c := range *courses {
		if len(c.EndAt) > 0 {
			endAt, err := time.Parse(time.RFC3339, c.EndAt)
			if err != nil {
				// we'll say that it's ended
				continue
			}

			if endAt.Before(time.Now()) {
				// course is over :(
				continue
			}
		}

		// this is here because there may not be a valid enrollment (say an observee is a TA)
		shouldAddCourse := false

		var validEnrollments []canvasEnrollment

		for _, e := range c.Enrollments {
			uID := e.UserID

			if e.Type != canvasEnrollmentTypeStudentEnrollment {
				// only allow student enrollments in gradedUsers
				continue
			}

			// we are now sure there is a valid enrollment for this course
			shouldAddCourse = true

			gradedUsers[c.ID] = append(gradedUsers[c.ID], uID)
			validEnrollments = append(validEnrollments, e)
		}

		c.Enrollments = validEnrollments

		if shouldAddCourse {
			validCourses = append(validCourses, c)
		}
	}

	return gradedUsers, &validCourses
}

// processOutcomeResults turns an unsorted list of canvas outcome alignments into a
// map like map[userID<uint64>]map[outcomeID<uint64>]canvasOutcomeAlignment
func processOutcomeResults(results *[]canvasOutcomeResult) (
	*map[uint64]map[uint64][]canvasOutcomeResult,
	error,
) {
	rs := map[uint64]map[uint64][]canvasOutcomeResult{}

	for _, r := range *results {
		userID, err := strconv.Atoi(r.Links.User)
		if err != nil {
			return nil, fmt.Errorf("error converting outcome result linked user ID (%s) to an integer: %w", r.Links.User, err)
		}
		uID := uint64(userID)

		outcomeID, err := strconv.Atoi(r.Links.LearningOutcome)
		if err != nil {
			return nil, fmt.Errorf(
				"error converting outcome result linked outcome ID (%s) to an integer: %w",
				r.Links.LearningOutcome,
				err)
		}
		oID := uint64(outcomeID)

		if rs[uID] == nil {
			rs[uID] = make(map[uint64][]canvasOutcomeResult)
		}

		rs[uID][oID] = append(rs[uID][oID], r)
	}

	return &rs, nil
}

func rdFromCanvasUserID(cuID uint64) (requestDetails, error) {
	tokens, err := canvas_tokens.List(db, &canvas_tokens.ListRequest{
		CanvasUserID: cuID,
		OrderBys:     []string{"inserted_at DESC"},
		Limit:        1,
	})
	if err != nil {
		return requestDetails{}, fmt.Errorf("error listing canvas tokens: %w", err)
	}

	if len(*tokens) < 1 {
		return requestDetails{}, nil
	}

	token := (*tokens)[0]

	return requestDetails{
		TokenID:      token.ID,
		Token:        token.Token,
		RefreshToken: token.RefreshToken,
	}, nil
}

func rdFromUserID(uID uint64) (requestDetails, error) {
	tokens, err := canvas_tokens.List(db, &canvas_tokens.ListRequest{
		UserID:   uID,
		OrderBys: []string{"inserted_at DESC"},
		Limit:    1,
	})
	if err != nil {
		return requestDetails{}, fmt.Errorf("error listing canvas tokens: %w", err)
	}

	if len(*tokens) < 1 {
		return requestDetails{}, nil
	}

	token := (*tokens)[0]

	return requestDetails{
		TokenID:      token.ID,
		Token:        token.Token,
		RefreshToken: token.RefreshToken,
	}, nil
}

// rdFromToken gets a requestDetails object from a db canvas token
func rdFromToken(tok canvas_tokens.CanvasToken) requestDetails {
	return requestDetails{
		TokenID:      tok.ID,
		Token:        tok.Token,
		RefreshToken: tok.RefreshToken,
	}
}

// sendJSON sends a JSON response.
func sendJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		handleISE(w, fmt.Errorf("error encoding json to send from gradesapi: %w", err))
		return
	}
}

// intFromQuery pulls a uint64 from a url.Values and handles it.
// If the return value is zero, return your request handler.
func intFromQuery(w http.ResponseWriter, name string, q url.Values) uint64 {
	n := q.Get(name)
	if len(n) < 1 || !util.ValidateIntegerString(n) {
		util.SendBadRequest(w, "missing or invalid "+name+" as query param")
		return 0
	}

	in, err := strconv.Atoi(n)
	if err != nil {
		util.SendBadRequest(w, "invalid "+name+" as query param")
	}

	return uint64(in)
}
