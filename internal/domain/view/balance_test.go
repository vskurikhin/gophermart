/*
 * This file was last modified at 2024-04-22 10:44 by Victor N. Skurikhin.
 * balance_test.go
 * $Id$
 */

package view

import (
	"github.com/stretchr/testify/assert"
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

				b := Balance{}
				assert.NotNil(t, b)
				assert.Equal(t, utils.BigFloatZero(), b.Current())
				assert.Equal(t, utils.BigFloatZero(), b.Sum())

				columns := []string{"login", "current", "withdrawn", "created_at", "update_at", "sum"}
				rows := utils.PgxRowsNext(
					columns, "test", utils.StringZero, utils.StringZero,
					utils.TimeZero(), utils.SQLNullTimeZero(), utils.SQLNullStringWith0(),
				)
				_, _, _, _, _, _, err := extractBalanceWithdrawn(rows)
				return b, err
			},
		},
		{
			name: "Test positive #2 Balance",
			get: func() (interface{}, error) {
				columns := []string{"login", "current", "withdrawn", "created_at", "update_at", "sum"}
				rows := utils.PgxRowsNext(
					columns, "test", utils.StringZero, utils.StringZero,
					utils.TimeZero(), utils.SQLNullTimeZero(), utils.SQLNullStringNull(),
				)
				_, _, _, _, _, _, err := extractBalanceWithdrawn(rows)
				return "ok", err
			},
		},
		{
			name: "Test positive #3 Balance",
			get: func() (interface{}, error) {
				columns := []string{"login", "current", "withdrawn", "created_at", "update_at", "sum"}
				rows := utils.PgxRowsNext(
					columns, "test", utils.StringZero, utils.StringZero,
					utils.TimeZero(), utils.SQLNullTimeZero(), utils.SQLNullStringZero(),
				)
				_, _, _, _, _, _, err := extractBalanceWithdrawn(rows)
				return "ok", err
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
				_, _, _, _, _, _, err := extractBalanceWithdrawn(rows)
				return nil, err
			},
		},
		{
			name: "Test negative #2 Balance",
			get: func() (interface{}, error) {
				columns := []string{"login", "current", "withdrawn", "created_at", "update_at", "sum"}
				rows := utils.PgxRowsNext(
					columns, "test", "", utils.StringZero,
					utils.TimeZero(), utils.SQLNullTimeZero(), utils.SQLNullStringWith0(),
				)
				_, _, _, _, _, _, err := extractBalanceWithdrawn(rows)
				return nil, err
			},
		},
		{
			name: "Test negative #3 Balance",
			get: func() (interface{}, error) {
				columns := []string{"login", "current", "withdrawn", "created_at", "update_at", "sum"}
				rows := utils.PgxRowsNext(
					columns, "test", utils.StringZero, "",
					utils.TimeZero(), utils.SQLNullTimeZero(), utils.SQLNullStringWith0(),
				)
				_, _, _, _, _, _, err := extractBalanceWithdrawn(rows)
				return nil, err
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
