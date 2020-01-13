package orautil

import (
	"reflect"
	"testing"
)

func TestBuildOracleUpsert(t *testing.T) {
	type args struct {
		tableName string
		checks    []UpsertableCheck
		data      []UpsertableColumn
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   []interface{}
		wantErr bool
	}{
		{
			name: "insert_and_update",
			args: struct {
				tableName string
				checks    []UpsertableCheck
				data      []UpsertableColumn
			}{tableName: "testtableA", checks: []UpsertableCheck{
				{
					Name:  "checkCol",
					Value: 1234,
				},
			}, data: []UpsertableColumn{
				{
					Name:       "columnA",
					Value:      "abcd",
					InsertOnly: true,
				},
				{
					Name:       "columnB",
					Value:      "hello, world",
					InsertOnly: false,
				},
			}},
			want: "MERGE INTO testtableA USING DUAL ON (checkCol=:1) WHEN NOT MATCHED THEN INSERT (columnA,columnB) VALUES (:2,:3) WHEN MATCHED THEN UPDATE SET columnB=:4",
			want1: []interface{}{
				1234,
				"abcd",
				"hello, world",
				"hello, world",
			},
			wantErr: false,
		},
		{
			name: "insert_only",
			args: struct {
				tableName string
				checks    []UpsertableCheck
				data      []UpsertableColumn
			}{tableName: "testtableA", checks: []UpsertableCheck{
				{
					Name:  "checkCol",
					Value: 1234,
				},
			}, data: []UpsertableColumn{
				{
					Name:       "columnA",
					Value:      "abcd",
					InsertOnly: true,
				},
				{
					Name:       "columnB",
					Value:      "hello, world",
					InsertOnly: true,
				},
			}},
			want: "MERGE INTO testtableA USING DUAL ON (checkCol=:1) WHEN NOT MATCHED THEN INSERT (columnA,columnB) VALUES (:2,:3)",
			want1: []interface{}{
				1234,
				"abcd",
				"hello, world",
			},
			wantErr: false,
		},
		{
			name: "multiple_checks_or",
			args: struct {
				tableName string
				checks    []UpsertableCheck
				data      []UpsertableColumn
			}{tableName: "testtableA", checks: []UpsertableCheck{
				{
					Name:  "checkCol",
					Value: 1234,
					Or:    true,
				},
				{
					Name:  "checkCol2",
					Value: 5678,
					// or is redundant here because this is the last check
					// however it's useful for testing-- shouldn't show up
					Or: true,
				},
			}, data: []UpsertableColumn{
				{
					Name:       "columnA",
					Value:      "abcd",
					InsertOnly: true,
				},
				{
					Name:       "columnB",
					Value:      "hello, world",
					InsertOnly: false,
				},
			}},
			want: "MERGE INTO testtableA USING DUAL ON (checkCol=:1 OR checkCol2=:2) WHEN NOT MATCHED THEN INSERT (columnA,columnB) VALUES (:3,:4) WHEN MATCHED THEN UPDATE SET columnB=:5",
			want1: []interface{}{
				1234,
				5678,
				"abcd",
				"hello, world",
				"hello, world",
			},
			wantErr: false,
		},
		{
			name: "multiple_checks_and",
			args: struct {
				tableName string
				checks    []UpsertableCheck
				data      []UpsertableColumn
			}{tableName: "testtableA", checks: []UpsertableCheck{
				{
					Name:  "checkCol",
					Value: 1234,
					Or:    false,
				},
				{
					Name:  "checkCol2",
					Value: 5678,
					// leaving Or blank-- should see AND anyway
				},
				{
					Name:  "checkCol3",
					Value: 9012,
					// or is redundant here because this is the last check
					// however it's useful for testing-- shouldn't show up
					Or: false,
				},
			}, data: []UpsertableColumn{
				{
					Name:       "columnA",
					Value:      "abcd",
					InsertOnly: true,
				},
				{
					Name:       "columnB",
					Value:      "hello, world",
					InsertOnly: false,
				},
			}},
			want: "MERGE INTO testtableA USING DUAL ON (checkCol=:1 AND checkCol2=:2 AND checkCol3=:3) WHEN NOT MATCHED THEN INSERT (columnA,columnB) VALUES (:4,:5) WHEN MATCHED THEN UPDATE SET columnB=:6",
			want1: []interface{}{
				1234,
				5678,
				9012,
				"abcd",
				"hello, world",
				"hello, world",
			},
			wantErr: false,
		},
		{
			name: "multiple_checks_and_or",
			args: struct {
				tableName string
				checks    []UpsertableCheck
				data      []UpsertableColumn
			}{tableName: "testtableA", checks: []UpsertableCheck{
				{
					Name:  "checkCol",
					Value: 1234,
					Or:    false,
				},
				{
					Name:  "checkCol2",
					Value: 5678,
					Or:    true,
				},
				{
					Name:  "checkCol3",
					Value: 9012,
					// or is redundant here because this is the last check
					// however it's useful for testing-- shouldn't show up
					Or: false,
				},
			}, data: []UpsertableColumn{
				{
					Name:       "columnA",
					Value:      "abcd",
					InsertOnly: true,
				},
				{
					Name:       "columnB",
					Value:      "hello, world",
					InsertOnly: false,
				},
			}},
			want: "MERGE INTO testtableA USING DUAL ON (checkCol=:1 AND checkCol2=:2 OR checkCol3=:3) WHEN NOT MATCHED THEN INSERT (columnA,columnB) VALUES (:4,:5) WHEN MATCHED THEN UPDATE SET columnB=:6",
			want1: []interface{}{
				1234,
				5678,
				9012,
				"abcd",
				"hello, world",
				"hello, world",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := BuildUpsert(tt.args.tableName, tt.args.checks, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildOracleUpsert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BuildOracleUpsert() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("BuildOracleUpsert() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
