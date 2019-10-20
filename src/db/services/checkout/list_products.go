package checkout

import (
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
)

type Product struct {
	ID        uint64
	Name      string
	ShortName string
	Price     float64
	Type      string
}

func GetAllProducts(db services.DB) (*[]Product, error) {
	query, args, err := util.Sq.
		Select("id", "name", "short_name", "price", "type").
		From("products").
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "error building list products sql")
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "error executing list products sql")
	}

	defer rows.Close()

	var ps []Product

	for rows.Next() {
		var p Product

		err := rows.Scan(&p.ID, &p.Name, &p.ShortName, &p.Price, &p.Type)
		if err != nil {
			return nil, errors.Wrap(err, "error scanning list products rows")
		}

		ps = append(ps, p)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "rows error in list products")
	}

	return &ps, nil
}
