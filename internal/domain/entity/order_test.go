/*
 * This file was last modified at 2024-04-21 18:18 by Victor N. Skurikhin.
 * order_test.go
 * $Id$
 */

package entity

import (
	"github.com/stretchr/testify/assert"
	"github.com/vskurikhin/gophermart/internal/storage"
	"github.com/vskurikhin/gophermart/internal/utils"
	"testing"
)

var orderColumns = []string{"login", "number", "status_id", "accrual", "uploaded_at", "created_at", "update_at"}

//goland:noinspection GoImportUsedAsName,GoUnhandledErrorResult
func TestOrderPositive(t *testing.T) {

	var tests = []struct {
		name string
		get  func() (interface{}, error)
	}{
		{
			name: "Test positive #1 Order",
			get: func() (interface{}, error) {

				o := NewOrder("test", utils.StringZero)
				assert.NotNil(t, o)
				assert.Equal(t, "test", o.Login())
				assert.Nil(t, o.Accrual())

				f := utils.BigFloatWith0()
				o.SetAccrual(&f)
				assert.Equal(t, utils.BigFloatWith0(), *o.Accrual())

				a := make(storage.TxArgs, 0)
				assert.True(t, len(a) == 0)
				a = o.AppendAccrualTo(a)
				assert.True(t, len(a) == 1)

				a = o.AppendSetStatusTo(a, utils.StringZero)
				assert.True(t, len(a) == 2)

				rows := utils.PgxRowsNext(
					orderColumns, "test", utils.StringZero, 0, utils.SqlNullStringWith0(),
					utils.SqlNullTimeZero(), utils.TimeZero(), utils.SqlNullTimeZero(),
				)
				return extractOrder(rows)
			},
		},
		{
			name: "Test positive #2 Order",
			get: func() (interface{}, error) {
				rows := utils.PgxRowsNext(
					orderColumns, "test", "", 0, utils.SqlNullStringWith0(),
					utils.SqlNullTimeZero(), utils.TimeZero(), utils.SqlNullTimeZero(),
				)
				return extractOrder(rows)
			},
		},
		{
			name: "Test positive #3 Order",
			get: func() (interface{}, error) {
				rows := utils.PgxRowsNext(
					orderColumns, "test", utils.StringZero, 0, utils.SqlNullStringZero(),
					utils.SqlNullTimeZero(), utils.TimeZero(), utils.SqlNullTimeZero(),
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
