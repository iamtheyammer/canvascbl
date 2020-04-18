package enrollments

import (
	"fmt"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
)

// UpsertRequest represents all the needed information to upsert an enrollment.
type UpsertRequest struct {
	//CanvasID               uint64
	CourseID               uint64
	UserCanvasID           uint64
	AssociatedUserCanvasID uint64
	EnrollmentRole         string
	EnrollmentState        string
	//CreatedAt              string
	//UpdatedAt              string
}

// Upsert upserts an enrollment.
func Upsert(db services.DB, req *[]UpsertRequest) error {
	q := util.Sq.
		Insert("enrollments").
		Columns(
			//"canvas_id",
			"course_id",
			"user_canvas_id",
			"associated_user_canvas_id",
			"enrollment_role",
			"enrollment_state",
			//"created_at",
			//"updated_at",
		).
		Suffix(
			"ON CONFLICT (course_id, user_canvas_id, associated_user_canvas_id) DO UPDATE SET " +
				"enrollment_role = EXCLUDED.enrollment_role, " +
				"enrollment_state = EXCLUDED.enrollment_state",
			//"updated_at = EXCLUDED.updated_at",
		)

	for _, r := range *req {
		var associatedUserCanvasID interface{}
		if r.AssociatedUserCanvasID > 0 {
			associatedUserCanvasID = r.AssociatedUserCanvasID
		}

		q = q.Values(
			//r.CanvasID,
			r.CourseID,
			r.UserCanvasID,
			associatedUserCanvasID,
			r.EnrollmentRole,
			r.EnrollmentState,
			//r.CreatedAt,
			//r.UpdatedAt,
		)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return fmt.Errorf("error building upsert enrollments sql: %w", err)
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("error executing upsert enrollments sql: %w", err)
	}

	return nil
}
