/*
 * This file was last modified at 2024-04-15 12:13 by Victor N. Skurikhin.
 * pgs_test.go
 * $Id$
 */

package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

//goland:noinspection GoImportUsedAsName,GoUnhandledErrorResult
func TestPgsStoragePositive(t *testing.T) {

	dbPool := DBPool()
	pgxPool, ok := dbPool.DBPool()
	assert.True(t, ok)
	storage := NewPgsStorage(pgxPool)

	var tests = []struct {
		name string
		sql  string
		args []interface{}
	}{
		{"Test positive #1", "SELECT version()", []interface{}{}},
		{"Test positive #2", "SELECT $1 + $2", []interface{}{1, 2}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := storage.sqlRow(test.name, test.sql, test.args...)
			assert.Nil(t, err)
		})
	}
}

/*
//goland:noinspection GoImportUsedAsName,GoUnhandledErrorResult
func TestNegative(t *testing.T) {

	var tests = []struct {
		name string
		test func()
	}{
		{
			name: "Test negative #1",
			test: func() {
				d := new(dbPoolHealth)
				d.log = logger.Get()

				defer func() {
					if r := recover(); r != nil {
						assert.NotNil(t, r)
					}
				}()
				d.pool = newPgxPool("user=jack password=secret host=pg.example.com port=5432 dbname=mydb sslmode=verify-ca pool_max_conns=nil")

			},
		},
		{
			name: "Test negative #2",
			test: func() {
				d := new(dbPoolHealth)
				d.log = logger.Get()

				defer func() {
					if r := recover(); r != nil {
						assert.NotNil(t, r)
					}
				}()
				d.pool = newPgxPool("")

			},
		},
		{
			name: "Test negative #3",
			test: func() {
				d := new(dbPoolHealth)
				d.dbPing()

				defer func() {
					if r := recover(); r != nil {
						assert.NotNil(t, r)
					}
				}()
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.test()
		})
	}
}
*/
