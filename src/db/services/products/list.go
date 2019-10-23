package products

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
)

type Product struct {
	ID        uint64
	StripeID  string
	Name      string
	ShortName string
	Price     float64
	Type      string
}

type ListRequest struct {
	ID        uint64
	StripeID  string
	ShortName string
}

// ListProducts lists all products
func ListProducts(db services.DB) (*[]Product, error) {
	query, args, err := util.Sq.
		Select("id", "stripe_id", "name", "short_name", "price", "type").
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

		err := rows.Scan(&p.ID, &p.StripeID, &p.Name, &p.ShortName, &p.Price, &p.Type)
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

// ListProduct lists a single product based on short name, id or both.
func ListProduct(db services.DB, req *ListRequest) (*Product, error) {
	if req.ID == 0 && len(req.StripeID) == 0 && len(req.ShortName) == 0 {
		return nil, nil
	}

	q := util.Sq.
		Select("id", "stripe_id", "name", "short_name", "price", "type").
		From("products")

	if req.ID != 0 {
		q = q.Where(sq.Eq{"id": req.ID})
	}

	if len(req.StripeID) > 0 {
		q = q.Where(sq.Eq{"stripe_id": req.StripeID})
	}

	if len(req.ShortName) > 0 {
		q = q.Where(sq.Eq{"short_name": req.ShortName})
	}

	query, args, err := q.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "error building list product sql")
	}

	var p Product

	row := db.QueryRow(query, args...)
	err = row.Scan(&p.ID, &p.StripeID, &p.Name, &p.ShortName, &p.Price, &p.Type)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrap(err, "error scanning list product row")
	}

	return &p, nil
}
