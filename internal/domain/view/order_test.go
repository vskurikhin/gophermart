/*
 * This file was last modified at 2024-04-21 18:48 by Victor N. Skurikhin.
 * order_test.go
 * $Id$
 */

package view

import (
	"github.com/stretchr/testify/assert"
	"github.com/vskurikhin/gophermart/internal/utils"
	"testing"
)

//goland:noinspection GoImportUsedAsName,GoUnhandledErrorResult
func TestOrderPositive(t *testing.T) {

	var tests = []struct {
		name string
		get  func() (interface{}, error)
	}{
		{
			name: "Test positive #1 Order",
			get: func() (interface{}, error) {

				b := Order{}
				assert.NotNil(t, b)
				assert.Equal(t, "", b.Number())
				assert.Equal(t, "", b.Status())
				assert.Nil(t, b.Accrual())

				columns := []string{"login", "number", "status_id", "accrual", "uploaded_at", "created_at", "update_at", "status"}
				rows := utils.PgxRowsNext(
					columns, "test", utils.StringZero, 0, utils.SqlNullStringWith0(),
					utils.SqlNullTimeZero(), utils.TimeZero(), utils.SqlNullTimeZero(), "NEW",
				)
				return extractOrder(rows)
			},
		},
		{
			name: "Test positive #1 Order",
			get: func() (interface{}, error) {

				b := Order{}
				assert.NotNil(t, b)
				assert.Equal(t, "", b.Number())
				assert.Equal(t, "", b.Status())
				assert.Nil(t, b.Accrual())

				columns := []string{"login", "number", "status_id", "accrual", "uploaded_at", "created_at", "update_at", "status"}
				rows := utils.PgxRowsNext(
					columns, "test", utils.StringZero, 0, utils.SqlNullStringZero(),
					utils.SqlNullTimeZero(), utils.TimeZero(), utils.SqlNullTimeZero(), "NEW",
				)
				return extractOrder(rows)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.get()
			assert.Nil(t, err)
			assert.NotNil(t, got)
		})
	}
}

//goland:noinspection GoImportUsedAsName,GoUnhandledErrorResult
func TestOrderNegative(t *testing.T) {

	var tests = []struct {
		name string
		get  func() (interface{}, error)
	}{
		{
			name: "Test negative #1 Order",
			get: func() (interface{}, error) {
				rows := &utils.TestRows{}
				return extractOrder(rows)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.get()
			assert.NotNil(t, err)
			assert.Nil(t, got)
		})
	}
}
