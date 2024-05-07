/*
 * This file was last modified at 2024-04-25 21:40 by Victor N. Skurikhin.
 * balance_test.go
 * $Id$
 */

package entity

import (
	"github.com/stretchr/testify/assert"
	"github.com/vskurikhin/gophermart/internal/storage"
	"github.com/vskurikhin/gophermart/internal/utils"
	"testing"
)

//goland:noinspection GoImportUsedAsName,GoUnhandledErrorResult
func TestBalancePositive(t *testing.T) {

	var tests = []struct {
		name string
		get  func() (interface{}, error)
	}{
		{
			name: "Test positive #1 Balance",
			get: func() (interface{}, error) {

				b := NewBalance("test", utils.BigFloatWith0())
				assert.NotNil(t, b)
				assert.Equal(t, utils.BigFloatWith0(), b.Current())

				a := make(storage.TxArgs, 0)
				assert.True(t, len(a) == 0)

				a = b.AppendInsertTo(a)
				assert.True(t, len(a) == 1)
				a = b.AppendAccrualTo(a, nil)
				assert.True(t, len(a) == 1)

				zero := utils.BigFloatWith0()
				a = b.AppendAccrualTo(a, &zero)
				assert.True(t, len(a) == 2)
				a = b.AppendWithdrawTo(a, nil)
				assert.True(t, len(a) == 2)

				a = b.AppendWithdrawTo(a, &zero)
				assert.True(t, len(a) == 3)

				columns := []string{"login", "current", "withdrawn", "created_at", "update_at"}
				rows := utils.PgxRowsNext(
					columns, "test", utils.StringZero, utils.StringZero, utils.TimeZero(), utils.SQLNullTimeZero(),
				)
				return extractBalance(rows)
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
func TestBalanceNegative(t *testing.T) {

	var tests = []struct {
		name string
		get  func() (interface{}, error)
	}{
		{
			name: "Test negative #1 Balance",
			get: func() (interface{}, error) {
				rows := &utils.TestRows{}
				_, err := extractBalance(rows)
				return nil, err
			},
		},
		{
			name: "Test negative #2 Balance",
			get: func() (interface{}, error) {
				columns := []string{"login", "current", "withdrawn", "created_at", "update_at"}
				rows := utils.PgxRowsNext(
					columns, "test", "", utils.StringZero, utils.TimeZero(), utils.SQLNullTimeNull(),
				)
				return extractBalance(rows)
			},
		},
		{
			name: "Test negative #3 Balance",
			get: func() (interface{}, error) {
				columns := []string{"login", "current", "withdrawn", "created_at", "update_at"}
				rows := utils.PgxRowsNext(
					columns, "test", utils.StringZero, "", utils.TimeZero(), utils.SQLNullTimeNull(),
				)
				return extractBalance(rows)
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
