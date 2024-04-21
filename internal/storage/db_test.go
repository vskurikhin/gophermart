/*
 * This file was last modified at 2024-04-15 21:30 by Victor N. Skurikhin.
 * db_test.go
 * $Id$
 */

package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

//goland:noinspection GoImportUsedAsName,GoUnhandledErrorResult
func TestDBPoolPositive(t *testing.T) {

	dbPool := GetDB()

	var tests = []struct {
		name string
	}{
		{"Test positive #1"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, dbPoolHealthInstance, dbPool)
			time.Sleep(3 * time.Second)
			_, ok := dbPool.DBPool()
			assert.True(t, ok)
		})
	}
}

//goland:noinspection GoImportUsedAsName,GoUnhandledErrorResult
func TestDBPoolNegative(t *testing.T) {

	var tests = []struct {
		name string
		test func()
	}{
		{
			name: "Test negative #1",
			test: func() {
				defer func() {
					if r := recover(); r != nil {
						assert.NotNil(t, r)
					}
				}()
				newPgxPool("")
			},
		},
		{
			name: "Test negative #2",
			test: func() {
				defer func() {
					if r := recover(); r != nil {
						assert.NotNil(t, r)
					}
				}()
				newPgxPool("user=jack password=secret host=pg.example.com port=5432 dbname=mydb sslmode=verify-ca pool_max_conns=nil")
			},
		},
		{
			name: "Test negative #3",
			test: func() {
				defer func() {
					if r := recover(); r != nil {
						assert.NotNil(t, r)
					}
				}()
				d := new(dbPoolHealth)
				d.dbPing()
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.test()
		})
	}
}
