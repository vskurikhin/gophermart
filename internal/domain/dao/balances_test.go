/*
 * This file was last modified at 2024-04-21 14:19 by Victor N. Skurikhin.
 * balances_test.go
 * $Id$
 */

package dao

import (
	"github.com/stretchr/testify/assert"
	"github.com/vskurikhin/gophermart/internal/utils"
	"go.uber.org/mock/gomock"
	"testing"
)

//goland:noinspection GoImportUsedAsName,GoUnhandledErrorResult
func TestDaoBalances(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var tests = []struct {
		name string
		get  func() (interface{}, error)
	}{
		{
			name: "Test #1 Balances.GetBalance",
			get: func() (interface{}, error) {
				columns := []string{"login", "current", "withdrawn", "created_at", "update_at"}
				m := MockStorageGetByString(
					ctrl, columns, "test", utils.StringZero, utils.StringZero, utils.TimeZero(), utils.SqlNullTimeNull(),
				)
				db := Balances(m)
				return db.GetBalance("test")
			},
		},
		{
			name: "Test #1 Balances.GetBalanceWithdraw",
			get: func() (interface{}, error) {
				columns := []string{"login", "current", "withdrawn", "created_at", "update_at", "sum"}
				m := MockStorageGetByString(
					ctrl, columns, "test", utils.StringZero, utils.StringZero,
					utils.TimeZero(), utils.SqlNullTimeNull(), utils.SqlNullStringNull(),
				)
				db := Balances(m)
				return db.GetBalanceWithdraw("test")
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
