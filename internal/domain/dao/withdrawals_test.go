/*
 * This file was last modified at 2024-04-22 10:44 by Victor N. Skurikhin.
 * withdrawals_test.go
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
func TestDaoWithdrawals(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var tests = []struct {
		name string
		get  func() (interface{}, error)
	}{
		{
			name: "Test #1 Withdrawals.GetAllWithdrawalsByLogin",
			get: func() (interface{}, error) {

				columns := []string{"login", "number", "sum", "status_id", "processed_at", "created_at", "update_at"}
				m := MockStorageGetAllForString(
					ctrl, columns, "test", utils.StringZero, utils.StringZero, 0,
					utils.SQLNullTimeNull(), utils.TimeZero(), utils.SQLNullTimeNull(),
				)
				dw := Withdrawals(m)
				dw.GetAllWithdrawalsByLogin("test")
				return "todo", nil
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
