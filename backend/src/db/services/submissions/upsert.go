package submissions

import (
	"fmt"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"time"
)

// Submission represents the metadata for a Canvas Submission.
type Submission struct {
	ID           uint64
	CanvasID     uint64
	CourseID     uint64
	AssignmentID uint64
	UserCanvasID uint64
	Attempt      uint64
	Score        float64
	// GraderID represents the ID of the user that graded it. If it was automatically graded, this value is negative.
	GraderID         int
	GradedAt         *time.Time
	Type             string
	SubmittedAt      *time.Time
	HTMLURL          string
	Late             bool
	Excused          bool
	Missing          bool
	LatePolicyStatus string
	PointsDeducted   float64
	SecondsLate      uint64
	ExtraAttempts    uint64
	PostedAt         *time.Time
	InsertedAt       time.Time
}

// Attachment represents a submission attachment in Canvas.
type Attachment struct {
	ID           uint64
	CanvasID     uint64
	SubmissionID uint64
	DisplayName  string
	Filename     string
	ContentType  string
	URL          string
	Size         uint64
	CreatedAt    *time.Time
	InsertedAt   time.Time
}

// UpsertRequest represents all the data required to upsert a Submission.
type UpsertRequest struct {
	CanvasID     uint64
	CourseID     uint64
	AssignmentID uint64
	UserCanvasID uint64
	Attempt      uint64
	Score        float64
	// GraderID represents the ID of the user that graded it. If it was automatically graded, this value is negative.
	GraderID         int
	GradedAt         time.Time
	Type             string
	SubmittedAt      *time.Time
	HTMLURL          string
	Late             bool
	Excused          bool
	Missing          bool
	LatePolicyStatus string
	PointsDeducted   float64
	SecondsLate      uint64
	ExtraAttempts    uint64
	PostedAt         *time.Time
}

// AttachmentUpsertRequest represents all the data required to upsert an attachment.
type AttachmentUpsertRequest struct {
	CanvasID     uint64
	SubmissionID uint64
	DisplayName  string
	Filename     string
	ContentType  string
	URL          string
	Size         uint64
	CreatedAt    *time.Time
}

// UpsertChunkSize represents the number of size of each upsert chunk.
// If your number of upserts is less than UpsertChunkSize, chunking is not necessary.
var UpsertChunkSize = services.CalculateChunkSize(19)

// Upsert upserts Submissions.
func Upsert(db services.DB, req *[]UpsertRequest) error {
	q := util.Sq.
		Insert("submissions").
		Columns(
			"canvas_id",
			"course_id",
			"assignment_id",
			"user_canvas_id",
			"attempt",
			"score",
			"grader_id",
			"graded_at",
			"submission_type",
			"submitted_at",
			"html_url",
			"late",
			"excused",
			"missing",
			"late_policy_status",
			"points_deducted",
			"seconds_late",
			"extra_attempts",
			"posted_at",
		).
		Suffix("ON CONFLICT (canvas_id) DO UPDATE SET " +
			"attempt = EXCLUDED.attempt, " +
			"score = EXCLUDED.score, " +
			"grader_id = EXCLUDED.grader_id, " +
			"graded_at = EXCLUDED.graded_at, " +
			"submission_type = EXCLUDED.submission_type, " +
			"submitted_at = EXCLUDED.submitted_at, " +
			"late = EXCLUDED.late, " +
			"missing = EXCLUDED.missing, " +
			"late_policy_status = EXCLUDED.late_policy_status, " +
			"points_deducted = EXCLUDED.points_deducted, " +
			"seconds_late = EXCLUDED.seconds_late, " +
			"extra_attempts = EXCLUDED.extra_attempts, " +
			"posted_at = EXCLUDED.posted_at",
		)

	for _, r := range *req {
		var score, graderID, gradedAt, submissionType, submittedAt, htmlURL, latePolicyStatus, pointsDeducted, extraAttempts interface{}

		if r.Score != 0 {
			score = r.Score
		}

		if r.GraderID != 0 {
			graderID = r.GraderID
		}

		if !r.GradedAt.IsZero() {
			gradedAt = r.GradedAt
		}

		if len(r.Type) > 0 {
			submissionType = r.Type
		}

		if r.SubmittedAt != nil && !(*r.SubmittedAt).IsZero() {
			submittedAt = r.SubmittedAt
		}

		if len(r.HTMLURL) > 0 {
			htmlURL = r.HTMLURL
		}

		if len(r.LatePolicyStatus) > 0 {
			latePolicyStatus = r.LatePolicyStatus
		}

		if r.PointsDeducted != 0 {
			pointsDeducted = r.PointsDeducted
		}

		if r.ExtraAttempts != 0 {
			extraAttempts = r.ExtraAttempts
		}

		q = q.Values(
			r.CanvasID,
			r.CourseID,
			r.AssignmentID,
			r.UserCanvasID,
			r.Attempt,
			score,
			graderID,
			gradedAt,
			submissionType,
			submittedAt,
			htmlURL,
			r.Late,
			r.Excused,
			r.Missing,
			latePolicyStatus,
			pointsDeducted,
			r.SecondsLate,
			extraAttempts,
			r.PostedAt,
		)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return fmt.Errorf("error building insert submissions sql: %w", err)
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("error executing insert submissions sql: %w", err)
	}

	return nil
}

// AttachmentsUpsertChunkSize is the number of attachments that can go in a chunk.
var AttachmentsUpsertChunkSize = services.CalculateChunkSize(8)

func UpsertAttachments(db services.DB, req *[]AttachmentUpsertRequest) error {
	q := util.Sq.
		Insert("submission_attachments").
		Columns(
			"canvas_id",
			"submission_id",
			"display_name",
			"filename",
			"content_type",
			"url",
			"size",
			"created_at",
		).
		Suffix("ON CONFLICT (canvas_id) DO NOTHING")

	for _, r := range *req {
		var displayName, contentType, size interface{}

		if len(r.DisplayName) > 0 {
			displayName = r.DisplayName
		}

		if len(r.ContentType) > 0 {
			contentType = r.ContentType
		}

		if r.Size > 0 {
			size = r.Size
		}

		q = q.Values(
			r.CanvasID,
			r.SubmissionID,
			displayName,
			r.Filename,
			contentType,
			r.URL,
			size,
			r.CreatedAt,
		)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return fmt.Errorf("error building insert submission attachments sql: %w", err)
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("error executing insert submission attachments sql: %w", err)
	}

	return nil
}
