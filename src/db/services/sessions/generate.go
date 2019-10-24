package sessions

import (
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
)

// Generate generates a session for the given userID using the given database
func Generate(db services.DB, userID uint64) (*string, error) {
	query, args, err := util.Sq.
		Insert("sessions").
		SetMap(map[string]interface{}{
			"canvas_user_id": userID,
		}).
		Suffix("RETURNING session_string").
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "error building insert session sql")
	}

	res := db.QueryRow(query, args...)

	var ss string

	err = res.Scan(&ss)
	if err != nil {
		return nil, errors.Wrap(err, "error executing insert session sql")
	}

	return &ss, nil
}
