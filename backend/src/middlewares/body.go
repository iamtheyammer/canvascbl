package middlewares

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io"
)

func DecodeJSONBody(body io.Reader, dest interface{}) error {
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&dest)
	if err != nil {
		return errors.Wrap(err, "error decoding request body to json")
	}

	return nil
}
