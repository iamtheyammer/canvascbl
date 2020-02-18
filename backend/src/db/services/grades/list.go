package grades

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
	"time"
)

type ListRequest struct {
	UserCanvasIDs *[]uint64
	Before        *time.Time
	After         *time.Time
	CourseIDs     *[]uint64
	ManualFetch   *bool
}

type Grade struct {
	ID           uint64
	CourseID     uint64
	Grade        string
	UserCanvasID uint64
	ManualFetch  bool
	InsertedAt   time.Time
}

func List(db services.DB, req *ListRequest) (*[]Grade, error) {
	q := util.Sq.
		Select(
			"DISTINCT ON (grades.course_id, grades.user_canvas_id) id",
			"user_canvas_id",
			"course_id",
			"grade",
			"manual_fetch",
			"inserted_at",
		).
		From("grades").
		OrderBy("grades.course_id, grades.user_canvas_id, grades.inserted_at DESC").
		Where(sq.Eq{"user_canvas_id": req.UserCanvasIDs})

	if req.Before != nil {
		// using this weird workaround because it doesn't work any other way
		q = q.Where("inserted_at < to_timestamp(?)", req.Before.Unix())
	}

	if req.After != nil {
		q = q.Where(sq.Gt{"inserted_at": req.After})
	}

	if req.CourseIDs != nil {
		q = q.Where(sq.Eq{"course_id": req.CourseIDs})
	}

	if req.ManualFetch != nil {
		q = q.Where(sq.Eq{"manual_fetch": *req.ManualFetch})
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
			&g.UserCanvasID,
			&g.CourseID,
			&g.Grade,
			&g.ManualFetch,
			&g.InsertedAt,
		)

		if err != nil {
			return nil, errors.Wrap(err, "error scanning list grades sql")
		}

		grades = append(grades, g)
	}

	return &grades, nil
}
