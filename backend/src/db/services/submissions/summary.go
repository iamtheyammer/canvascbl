package submissions

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
)

/*
SubmissionSummary represents a summary of submissions.

Since zeroes are significant (it is legitimate to have zero submissions), nullable values are pointers.
*/
type SubmissionSummary struct {
	Unsubmitted       uint64
	LateUnsubmitted   *uint64
	PendingReview     uint64
	LatePendingReview *uint64
	Submitted         uint64
	LateSubmitted     *uint64
	Graded            uint64
	LateGraded        *uint64
}

// CourseUserSummaryRequest represents a request to GetCourseUserSummary.
type CourseUserSummaryRequest struct {
	CourseID uint64
	UserIDs  []uint64

	/*
		SeparateLate allows you to separate late submissions.

		This means that, for example, late graded assignments would not be included in Graded, but in LateGraded.
		To find the total number of graded assignments, you would take Graded+LateGraded.
	*/
	SeparateLate bool
}

// CourseUserSummaryResponse represents a response from GetCourseUserSummary.
type CourseUserSummaryResponse map[uint64]SubmissionSummary

type summaryRow struct {
	Count         uint64
	UserID        uint64
	WorkflowState WorkflowState
	Late          bool
}

func GetCourseUserSummary(db services.DB, req *CourseUserSummaryRequest) (*CourseUserSummaryResponse, error) {
	cols := []string{
		"COUNT(*) AS count",
		"user_canvas_id",
		"workflow_state",
	}

	groupBys := []string{
		"workflow_state",
		"user_canvas_id",
	}

	if req.SeparateLate {
		cols = append(cols, "late")
		groupBys = append(groupBys, "late")
	}

	q := util.Sq.
		Select(cols...).
		From("submissions").
		GroupBy(groupBys...).
		OrderBy("user_canvas_id")

	if req.CourseID > 0 {
		q = q.Where(sq.Eq{"course_id": req.CourseID})
	}

	if len(req.UserIDs) > 0 {
		q = q.Where(sq.Eq{"user_canvas_id": req.UserIDs})
	}

	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building get course user summary sql: %w", err)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing get course user summary sql: %w", err)
	}

	defer rows.Close()

	summary := CourseUserSummaryResponse{}
	for rows.Next() {
		var r summaryRow
		if req.SeparateLate {
			err = rows.Scan(
				&r.Count,
				&r.UserID,
				&r.WorkflowState,
				&r.Late,
			)
		} else {
			err = rows.Scan(
				&r.Count,
				&r.UserID,
				&r.WorkflowState,
			)
		}

		if err != nil {
			return nil, fmt.Errorf("error scanning get course user summary sql: %w", err)
		}

		s := summary[r.UserID]

		switch r.WorkflowState {
		case WorkflowStateGraded:
			if r.Late {
				s.LateGraded = &r.Count
			} else {
				s.Graded = r.Count
			}
		case WorkflowStatePendingReview:
			if r.Late {
				s.LatePendingReview = &r.Count
			} else {
				s.PendingReview = r.Count
			}
		case WorkflowStateSubmitted:
			if r.Late {
				s.LateSubmitted = &r.Count
			} else {
				s.Submitted = r.Count
			}
		case WorkflowStateUnsubmitted:
			if r.Late {
				s.LateUnsubmitted = &r.Count
			} else {
				s.Unsubmitted = r.Count
			}
		}

		summary[r.UserID] = s
	}

	return &summary, nil
}
