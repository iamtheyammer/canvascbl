package util

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	_ "github.com/godror/godror"
	"github.com/iamtheyammer/canvascbl/backend/src/env"
	"github.com/pkg/errors"
)

var (
	DB                = ConnectDB()
	placeholderFormat = sq.Colon
	Sq                = sq.StatementBuilder.PlaceholderFormat(placeholderFormat)
)

// ConnectDB connects to the database. Should only be called once.
func ConnectDB() *sql.DB {
	db, err := sql.Open("godror", env.DatabaseDSN)
	if err != nil {
		panic(errors.Wrap(err, "failed to call sql.open"))
	}

	err = db.Ping()
	if err != nil {
		panic(errors.Wrap(err, "failed to ping database"))
	}

	return db
}
