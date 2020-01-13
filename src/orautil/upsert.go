package orautil

import (
	"fmt"
	"github.com/pkg/errors"
)

type UpsertableColumn struct {
	// Name is the name of the column. Will not be parameterized as it should
	// not be user-submitted.
	Name  string
	Value interface{}
	// Whether this column should be used only for inserts (true) or for
	// updates too (false). Not used for checks.
	InsertOnly bool
}

type UpsertableCheck struct {
	// Name is the name of the column. Will not be parameterized as it should
	// not be user-submitted.
	Name string
	// Value represents a single value, not an array that we will check against.
	Value interface{}
	// Or is if you want this check to be 'or-ed' with the next check, giving
	// you something like col1=val1 OR col1=val2. Making this false (default)
	// will use AND.
	Or bool
}

var (
	NoChecksUpsertError = errors.New("no checks present in upsert")
	NoDataUpsertError   = errors.New("no data present in upsert")
)

// BuildOracleUpsert builds a working upsert in Oracle and returns a string of
// sql, an interface slice as args and an error if one occurred. No columns
// are required to be upserted-- sending no columns in the data param that have
// InsertOnly=true will simply result in the database doing nothing on a
// what-would-be duplicate insert, similar to a PostgreSQL ON CONFLICT DO
// NOTHING.
func BuildUpsert(tableName string, checks []UpsertableCheck, data []UpsertableColumn) (string, []interface{}, error) {
	if len(checks) < 1 {
		return "", []interface{}{}, NoChecksUpsertError
	}

	if len(data) < 1 {
		return "", []interface{}{}, NoDataUpsertError
	}

	// begin the query
	query := "MERGE INTO " + tableName + " USING DUAL ON ("

	// start to hold args
	// the length of args will be used to determine the 'arg number' (:n)
	var args []interface{}

	// first, handle checks

	// number of checks to figure out which is the last
	cl := len(checks) - 1
	for i, chk := range checks {
		// we append to args first so that we end up with :1 instead of :0
		args = append(args, chk.Value)

		// what to combine the statements with
		combinator := "AND"
		if chk.Or {
			combinator = "OR"
		}

		if i != cl {
			// there is another check
			query += fmt.Sprintf("%s=:%d %s ", chk.Name, len(args), combinator)
		} else {
			// this is the last check
			// combinators aren't used because there isn't another check
			query += fmt.Sprintf("%s=:%d) ", chk.Name, len(args))
		}
	}

	var updatableColumns []UpsertableColumn

	// second, handle insert case

	query += "WHEN NOT MATCHED THEN INSERT ("
	insertSuffix := "VALUES ("
	cl = len(data) - 1
	for i, ins := range data {
		query += ins.Name
		args = append(args, ins.Value)

		if i != cl {
			// there is another column
			query += ","
			insertSuffix += fmt.Sprintf(":%d,", len(args))
		} else {
			query += ") "
			insertSuffix += fmt.Sprintf(":%d)", len(args))
		}

		if !ins.InsertOnly {
			updatableColumns = append(updatableColumns, ins)
		}
	}
	query += insertSuffix

	// third, handle update case if it exists

	if len(updatableColumns) > 0 {
		// we add the space here instead of on the last append because we may
		// not have anything to update
		query += " WHEN MATCHED THEN UPDATE SET "
		cl := len(updatableColumns) - 1
		for i, upd := range updatableColumns {
			args = append(args, upd.Value)
			query += fmt.Sprintf("%s=:%d", upd.Name, len(args))

			if i != cl {
				// there is another column
				query += ","
			}
		}
	}

	return query, args, nil
}
