package util

import (
	"database/sql"
	"github.com/iamtheyammer/canvascbl/backend/src/env"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

var DB = ConnectDB()

// ConnectDB connects to the database. Should only be called once.
func ConnectDB() *sql.DB {
	db, err := sql.Open("postgres", env.DatabaseDSN)
	if err != nil {
		panic(errors.Wrap(err, "failed to call sql.open"))
	}

	err = db.Ping()
	if err != nil {
		panic(errors.Wrap(err, "failed to ping database"))
	}

	return db
}
