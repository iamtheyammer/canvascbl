package gradesapi

import (
	"errors"
	"fmt"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/canvas_tokens"
	"strconv"
	"time"
)

const spring20EnrollmentTermID = 4
const spring20DLEnrollmentTermID = 3

type processedCourses struct {
	// Current is all classes that are "in session"
	Current []canvasCourse
	// ByEnrollment is a map between a course ID and enrolled user IDs
	ByEnrollment map[uint64][]uint64
}

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

		for _, e := range c.Enrollments {
			uID := e.UserID

			switch e.Type {
			case canvasEnrollmentTypeObserverEnrollment:
				uID = e.AssociatedUserID
			case canvasEnrollmentTypeStudentEnrollment:
			default:
				continue
			}

			// we are now sure there is a valid enrollment for this course
			shouldAddCourse = true

			gradedUsers[c.ID] = append(gradedUsers[c.ID], uID)
		}

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
