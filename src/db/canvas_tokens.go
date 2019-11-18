package db

import (
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/canvas_tokens"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
)

func InsertCanvasToken(req *canvas_tokens.InsertRequest) error {
	err := canvas_tokens.Insert(util.DB, req)
	if err != nil {
		return errors.Wrap(err, "error inserting a canvas token")
	}

	return nil
}

func ListCanvasTokens(req *canvas_tokens.ListRequest) (*[]canvas_tokens.CanvasToken, error) {
	cts, err := canvas_tokens.List(util.DB, req)
	if err != nil {
		return nil, errors.Wrap(err, "error listing canvas tokens")
	}

	return cts, nil
}
