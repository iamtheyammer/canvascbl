package db

import (
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/google_users"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
)

func UpsertGoogleProfile(req *google_users.UpsertRequest) (*uint64, error) {
	id, err := google_users.Upsert(util.DB, req)
	if err != nil {
		return nil, errors.Wrap(err, "error upserting google profile")
	}
	return id, nil
}

func ListGoogleProfiles(req *google_users.ListRequest) (*[]google_users.GoogleUser, error) {
	gus, err := google_users.List(util.DB, req)
	if err != nil {
		return nil, errors.Wrap(err, "error listing google profiles")
	}

	return gus, nil
}

func UpdateGoogleProfile(req *google_users.UpdateRequest) error {
	err := google_users.Update(util.DB, req)
	if err != nil {
		return errors.Wrap(err, "error updating google profiles")
	}

	return nil
}
