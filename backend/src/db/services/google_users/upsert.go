package google_users

import (
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
)

// UpsertRequest represents a row. Pointer fields represent nullable values.
type UpsertRequest struct {
	// users.id
	UserID *uint64
	// id for the user from Google
	GoogleID string
	Email    string
	// first name
	GivenName *string
	// last name
	FamilyName *string
	// GivenName FamilyName
	Name *string
	// URL to profile picture for user
	ProfilePictureURL *string
	// user-customizable
	Gender *string
	// abbreviated as hd
	HostedDomain *string
}

func Upsert(db services.DB, req *UpsertRequest) (*uint64, error) {
	query, args, err := util.Sq.
		Insert("google_users").
		SetMap(map[string]interface{}{
			"users_id":            req.UserID,
			"google_id":           req.GoogleID,
			"email":               req.Email,
			"given_name":          req.GivenName,
			"family_name":         req.FamilyName,
			"name":                req.Name,
			"profile_picture_url": req.ProfilePictureURL,
			"gender":              req.Gender,
			"hd":                  req.HostedDomain,
		}).
		Suffix("ON CONFLICT ON CONSTRAINT google_users_google_id_key DO UPDATE SET " +
			"name = EXCLUDED.name, " +
			"given_name = EXCLUDED.given_name, " +
			"family_name = EXCLUDED.family_name, " +
			"gender = EXCLUDED.gender " +
			"RETURNING id").
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "error building insert google user sql")
	}

	row := db.QueryRow(query, args...)

	var id uint64
	err = row.Scan(&id)
	if err != nil {
		return nil, errors.Wrap(err, "error getting id from upsert google users")
	}

	return &id, nil
}
