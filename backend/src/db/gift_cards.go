package db

import (
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/gift_cards"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
)

func ListGiftCards(req *gift_cards.ListRequest) (*[]gift_cards.GiftCard, error) {
	gcs, err := gift_cards.List(util.DB, req)
	if err != nil {
		return nil, errors.Wrap(err, "error listing gift cards")
	}

	return gcs, nil
}

func UpdateGiftCards(req *gift_cards.UpdateRequest) error {
	err := gift_cards.Update(util.DB, req)
	if err != nil {
		return errors.Wrap(err, "error updating gift cards")
	}

	return nil
}

func InsertGiftCards(req *gift_cards.InsertRequest) (*[]gift_cards.GiftCard, error) {
	gcs, err := gift_cards.Insert(util.DB, req)
	if err != nil {
		return nil, errors.Wrap(err, "error inserting gift cards")
	}

	return gcs, nil
}
