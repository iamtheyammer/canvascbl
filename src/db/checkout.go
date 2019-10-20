package db

import (
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/checkout"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
)

func CheckoutListProducts() (*[]checkout.Product, error) {
	products, err := checkout.GetAllProducts(util.DB)
	if err != nil {
		return nil, errors.Wrap(err, "error listing products")
	}

	return products, nil
}
