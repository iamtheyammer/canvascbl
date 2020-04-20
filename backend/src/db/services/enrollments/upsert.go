package enrollments

import (
	"fmt"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
)

// UpsertRequest represents all the needed information to upsert an enrollment.
type UpsertRequest struct {
	CanvasID               uint64
	CourseID               uint64
	UserCanvasID           uint64
	AssociatedUserCanvasID uint64
	Type                   Type
	Role                   Role
	State                  string
	CreatedAt              string
	UpdatedAt              string
}

// Upsert upserts an enrollment.
func Upsert(db services.DB, req *[]UpsertRequest) error {
	q := util.Sq.
		Insert("enrollments").
		Columns(
			"canvas_id",
			"course_id",
			"user_canvas_id",
			"associated_user_canvas_id",
			"enrollment_type",
			"enrollment_role",
			"enrollment_state",
			"created_at",
			"updated_at",
		).
		Suffix(
			"ON CONFLICT (course_id, user_canvas_id) DO UPDATE SET " +
				"enrollment_type = EXCLUDED.enrollment_type, " +
				"enrollment_role = EXCLUDED.enrollment_role, " +
				"enrollment_state = EXCLUDED.enrollment_state, " +
				"updated_at = EXCLUDED.updated_at",
		)

	for _, r := range *req {
		// BUT WHY? Well, since an observer can have more than one enrollment and we don't get
		// enrollment IDs from courses, we must just ignore them :(
		// However, we already save observer-observee links in the observees table so we're ok.
		if r.Role == RoleObserver {
			continue
		}

		var (
			canvasID, associatedUserCanvasID, createdAt, updatedAt interface{}
		)

		if r.CanvasID > 0 {
			canvasID = r.CanvasID
		}

		if r.AssociatedUserCanvasID > 0 {
			associatedUserCanvasID = r.AssociatedUserCanvasID
		}

		if len(r.CreatedAt) > 0 {
			createdAt = r.CreatedAt
		}

		if len(r.UpdatedAt) > 0 {
			updatedAt = r.UpdatedAt
		}

		q = q.Values(
			canvasID,
			r.CourseID,
			r.UserCanvasID,
			associatedUserCanvasID,
			r.Type,
			r.Role,
			r.State,
			createdAt,
			updatedAt,
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
