package google_users

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
)

type GoogleUser struct {
	ID                uint64
	UsersID           uint64
	GoogleID          string
	Email             string
	Name              string
	GivenName         string
	FamilyName        string
	ProfilePictureURL string
	Gender            string
	HostedDomain      string
}

type ListRequest struct {
	ID           uint64
	UsersID      uint64
	GoogleID     string
	Email        string
	Name         string
	GivenName    string
	FamilyName   string
	Gender       string
	HostedDomain string

	Limit    uint64
	Offset   uint64
	OrderBys []string
}

func List(db services.DB, req *ListRequest) (*[]GoogleUser, error) {
	q := util.Sq.
		Select(
			"id",
			"users_id",
			"google_id",
			"email",
			"name",
			"given_name",
			"family_name",
			"profile_picture_url",
			"gender",
			"hd",
		).
		From("google_users")

	if req.ID > 0 {
		q = q.Where(sq.Eq{"id": req.ID})
	}

	if req.UsersID > 0 {
		q = q.Where(sq.Eq{"users_id": req.UsersID})
	}

	if len(req.GoogleID) > 0 {
		q = q.Where(sq.Eq{"google_id": req.GoogleID})
	}

	if len(req.Email) > 0 {
		q = q.Where(sq.Eq{"email": req.Email})
	}

	if len(req.Name) > 0 {
		q = q.Where(sq.Eq{"name": req.Name})
	}

	if len(req.GivenName) > 0 {
		q = q.Where(sq.Eq{"given_name": req.GivenName})
	}

	if len(req.FamilyName) > 0 {
		q = q.Where(sq.Eq{"family_name": req.FamilyName})
	}

	if len(req.Gender) > 0 {
		q = q.Where(sq.Eq{"gender": req.Gender})
	}

	if len(req.HostedDomain) > 0 {
		q = q.Where(sq.Eq{"hd": req.HostedDomain})
	}

	if req.Limit > 0 {
		q = q.Limit(req.Limit)
	} else {
		q = q.Limit(services.DefaultSelectLimit)
	}

	if req.Offset > 0 {
		q = q.Offset(req.Offset)
	}

	if len(req.OrderBys) > 0 {
		q = q.OrderBy(req.OrderBys...)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "error building list google profiles sql")
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "error querying for google profiles")
	}

	defer rows.Close()

	var gus []GoogleUser

	for rows.Next() {
		var (
			gu                                                         GoogleUser
			usersID                                                    sql.NullInt64
			name, givenName, familyName, profilePictureURL, gender, hd sql.NullString
		)

		err := rows.Scan(
			&gu.ID,
			&usersID,
			&gu.GoogleID,
			&gu.Email,
			&name,
			&givenName,
			&familyName,
			&profilePictureURL,
			&gender,
			&hd,
		)
		if err != nil {
			return nil, errors.Wrap(err, "error scanning google profiles")
		}

		if usersID.Valid {
			gu.UsersID = uint64(usersID.Int64)
		}

		if name.Valid {
			gu.Name = name.String
		}

		if givenName.Valid {
			gu.GivenName = givenName.String
		}

		if familyName.Valid {
			gu.FamilyName = familyName.String
		}

		if profilePictureURL.Valid {
			gu.ProfilePictureURL = profilePictureURL.String
		}

		if gender.Valid {
			gu.Gender = gender.String
		}

		if hd.Valid {
			gu.HostedDomain = hd.String
		}

		gus = append(gus, gu)
	}

	return &gus, nil
}
