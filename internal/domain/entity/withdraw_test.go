/*
 * This file was last modified at 2024-04-22 10:44 by Victor N. Skurikhin.
 * withdraw_test.go
 * $Id$
 */

package entity

import (
	"github.com/stretchr/testify/assert"
	"github.com/vskurikhin/gophermart/internal/storage"
	"github.com/vskurikhin/gophermart/internal/utils"
	"testing"
)

var withdrawColumns = []string{"login", "number", "sum", "status_id", "processed_at", "created_at", "update_at"}

//goland:noinspection GoImportUsedAsName,GoUnhandledErrorResult
func TestWithdrawPositive(t *testing.T) {

	var tests = []struct {
		name string
		get  func() (interface{}, error)
	}{
		{
			name: "Test positive #1 Withdraw",
			get: func() (interface{}, error) {

				u := NewWithdraw("test", utils.StringZero, utils.BigFloatWith0())
				assert.NotNil(t, u)
				assert.Equal(t, utils.StringZero, u.Number())
				assert.Equal(t, utils.BigFloatWith0(), u.Sum())
				assert.Nil(t, u.ProcessedAt())
				assert.Equal(t, utils.TimeEpoch(), u.CreatedAt())

				a := make(storage.TxArgs, 0)
				assert.True(t, len(a) == 0)

				a = u.AppendInsertTo(a)
				assert.True(t, len(a) == 1)

				rows := utils.PgxRowsNext(
					withdrawColumns, "test", utils.StringZero, utils.StringZero, 0,
					utils.SQLNullTimeZero(), utils.TimeZero(), utils.SQLNullTimeZero(),
				)
				return extractWithdraw(rows)
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
func TestWithdrawNegative(t *testing.T) {

	var tests = []struct {
		name string
		get  func() (interface{}, error)
	}{
		{
			name: "Test negative #1 Withdraw",
			get: func() (interface{}, error) {
				rows := &utils.TestRows{}
				return extractWithdraw(rows)
			},
		},
		{
			name: "Test negative #2 Withdraw",
			get: func() (interface{}, error) {
				rows := utils.PgxRowsNext(
					withdrawColumns, "test", utils.StringZero, "", 0,
					utils.SQLNullTimeZero(), utils.TimeZero(), utils.SQLNullTimeZero(),
				)
				return extractWithdraw(rows)
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
