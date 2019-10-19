package db

import (
	"github.com/iamtheyammer/canvascbl/backend/src/db/canvas/courses"
	coursessvc "github.com/iamtheyammer/canvascbl/backend/src/db/services/courses"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
)

// UpsertMultipleCourses takes a course response from Canvas and upserts the courses
func UpsertMultipleCourses(cr *string) {
	ccr, err := courses.FromJSON(cr)
	if err != nil {
		handleError(err)
		return
	}

	var urs []coursessvc.UpsertRequest

	for _, c := range *ccr {
		urs = append(urs, coursessvc.UpsertRequest{
			Name:       c.Name,
			CourseCode: c.CourseCode,
			State:      c.State,
			UUID:       c.UUID,
			CourseID:   c.ID,
		})
	}

	err = coursessvc.UpsertMultiple(util.DB, &urs)
	if err != nil {
		handleError(errors.Wrap(err, "error upserting multiple courses"))
		return
	}

	return
}
