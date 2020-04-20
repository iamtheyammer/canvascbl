package enrollments

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"time"
)

/*
Role represents an enrollment role.

Enrollment roles can be customized and created by Canvas administrators, so the best way to figure out a user's
actual role in a course is to use Type instead.

Constants are provided for the default roles in Canvas.
*/
type Role string

// ToType turns a Role into a Type. If a Role is not recognized, it will simply be typecasted to a Type.
func (r Role) ToType() Type {
	switch r {
	case RoleTeacher:
		return TypeTeacher
	case RoleTA:
		return TypeTA
	case RoleDesigner:
		return TypeDesigner
	case RoleStudent:
		return TypeStudent
	case RoleObserver:
		return TypeObserver
	default:
		return Type(r)
	}
}

// Type represents an enrollment type. Constants are provided for all possible enrollment types in Canvas.
type Type string

// State represents an enrollment state. Constants are provided for all possible enrollment states in Canvas.
type State string

const (
	// RoleTeacher represents the TeacherEnrollment enrollment role.
	RoleTeacher = "TeacherEnrollment"
	// RoleTA represents the TaEnrollment enrollment role.
	RoleTA = "TaEnrollment"
	// RoleDesigner represents the DesignerEnrollment role.
	RoleDesigner = "DesignerEnrollment"
	// RoleStudent represents the StudentEnrollment role.
	RoleStudent = "StudentEnrollment"
	// RoleObserver represents the ObserverEnrollment role.
	RoleObserver = "ObserverEnrollment"
	// TypeTeacher represents the enrollment type of teacher.
	TypeTeacher = "teacher"
	// TypeTA represents the enrollment type of ta.
	TypeTA = "ta"
	// TypeDesigner represents the enrollment type of designer.
	TypeDesigner = "designer"
	// TypeStudent represents the enrollment type of student.
	TypeStudent = "student"
	// TypeObserver represents the enrollment type of observer.
	TypeObserver = "observer"
	// StateActive represents an active enrollment.
	StateActive = "active"
	// StateInvitedOrPending represents an enrollment that is pending or that the user has been invited to.
	StateInvitedOrPending = "invited_or_pending"
	// StateCompleted represents a completed enrollment.
	StateCompleted = "completed"
)

// Enrollment represents an enrollment.
type Enrollment struct {
	ID           uint64
	CourseID     uint64
	UserCanvasID uint64
	Type         Type
	Role         Role
	State        State
	InsertedAt   time.Time
}

// ListRequest is the request for enrollments.List.
type ListRequest struct {
	ID           uint64
	CourseID     uint64
	CourseIDs    []uint64
	UserCanvasID uint64
	UserID       uint64
	Type         Type
	Role         Role
	State        State
}

// List lists enrollments.
func List(db services.DB, req *ListRequest) (*[]Enrollment, error) {
	q := util.Sq.
		Select(
			"enrollments.id",
			"enrollments.course_id",
			"enrollments.user_canvas_id",
			"enrollments.enrollment_type",
			"enrollments.enrollment_role",
			"enrollments.enrollment_state",
			"enrollments.inserted_at",
		).
		From("enrollments")

	if req.ID > 0 {
		q = q.Where(sq.Eq{"enrollments.id": req.ID})
	}

	if req.CourseID > 0 {
		q = q.Where(sq.Eq{"enrollments.course_id": req.CourseID})
	}

	if len(req.CourseIDs) > 0 {
		q = q.Where(sq.Eq{"enrollments.course_id": req.CourseIDs})
	}

	if req.UserCanvasID > 0 {
		q = q.Where(sq.Eq{"enrollments.user_canvas_id": req.UserCanvasID})
	}

	if req.UserID > 0 {
		q = q.
			Join("users ON enrollments.user_canvas_id = users.canvas_user_id").
			Where(sq.Eq{"users.id": req.UserID})
	}

	if len(req.Type) > 0 {
		q = q.Where(sq.Eq{"enrollments.enrollment_type": req.Type})
	}

	if len(req.Role) > 0 {
		q = q.Where(sq.Eq{"enrollments.enrollment_role": req.Role})
	}

	if len(req.State) > 0 {
		q = q.Where(sq.Eq{"enrollments.enrollment_state": req.State})
	}

	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building list enrollments sql: %w", err)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing list enrollments sql: %w", err)
	}

	defer rows.Close()

	var es []Enrollment
	for rows.Next() {
		var e Enrollment
		err := rows.Scan(
			&e.ID,
			&e.CourseID,
			&e.UserCanvasID,
			&e.Type,
			&e.Role,
			&e.State,
			&e.InsertedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning list enrollments sql: %w", err)
		}

		es = append(es, e)
	}

	return &es, nil
}
