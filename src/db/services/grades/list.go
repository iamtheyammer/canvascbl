package grades

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
	"time"
)

type ListRequest struct {
	UserLTIUserID *string
	Before        *time.Time
	After         *time.Time
	CourseIDs     *[]uint64
}

type Grade struct {
	ID            uint64
	CourseID      uint64
	Grade         string
	UserLTIUserID string
	InsertedAt    time.Time
}

func List(db services.DB, req *ListRequest) (*[]Grade, error) {
	q := util.Sq.
		Select("DISTINCT ON (grades.course_id) id", "course_id", "grade", "inserted_at").
		From("grades").
		Where(sq.Eq{"user_lti_user_id": req.UserLTIUserID})

	if req.Before != nil {
		q = q.Where(sq.Lt{"inserted_at": req.Before})
	}

	if req.After != nil {
		q = q.Where(sq.Gt{"inserted_at": req.After})
	}

	if req.CourseIDs != nil {
		q = q.Where(sq.Eq{"course_id": req.CourseIDs})
	}

	query, args, err := q.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "error building list grades sql")
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "error opening rows in list grades sql")
	}

	defer rows.Close()

	var grades []Grade

	for rows.Next() {
		var g Grade

		err := rows.Scan(
			&g.ID,
			&g.CourseID,
			&g.Grade,
			&g.InsertedAt,
		)

		if err != nil {
			return nil, errors.Wrap(err, "error scanning list grades sql")
		}

		grades = append(grades, g)
	}

	return &grades, nil
}
