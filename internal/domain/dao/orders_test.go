/*
 * This file was last modified at 2024-04-22 10:44 by Victor N. Skurikhin.
 * orders_test.go
 * $Id$
 */

package dao

import (
	"github.com/stretchr/testify/assert"
	"github.com/vskurikhin/gophermart/internal/domain/entity"
	"github.com/vskurikhin/gophermart/internal/utils"
	"go.uber.org/mock/gomock"
	"testing"
)

//goland:noinspection GoImportUsedAsName,GoUnhandledErrorResult
func TestDaoOrders(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var tests = []struct {
		name string
		get  func() (interface{}, error)
	}{
		{
			name: "Test #1 Orders.GetOrderByNumber",
			get: func() (interface{}, error) {
				columns := []string{"login", "number", "status_id", "accrual", "uploaded_at", "created_at", "update_at"}
				m := MockStorageGetByString(
					ctrl, columns, "test", utils.StringZero, 0, utils.SQLNullStringNull(),
					utils.SQLNullTimeNull(), utils.TimeZero(), utils.SQLNullTimeNull(),
				)
				do := Orders(m)
				return do.GetOrderByNumber("0")
			},
		},
		{
			name: "Test #2 GetBalanceWithdraw",
			get: func() (interface{}, error) {
				columns := []string{"login", "number", "accrual", "uploaded_at", "created_at", "update_at", "status"}
				m := MockStorageGetAllForString(
					ctrl, columns, "test", utils.StringZero, utils.SQLNullStringNull(), utils.SQLNullTimeNull(),
					utils.TimeZero(), utils.SQLNullTimeNull(), utils.SQLNullStringNull(),
				)
				do := Orders(m)
				do.GetAllOrdersForLogin("test")
				return "todo", nil
			},
		},
		{
			name: "Test #3 Insert",
			get: func() (interface{}, error) {
				columns := []string{"login", "number", "status_id", "accrual", "uploaded_at", "created_at", "update_at"}
				m := MockStorageSaveOrder(
					ctrl, columns, "test", utils.StringZero, 0,
					utils.SQLNullStringNull(), utils.SQLNullTimeNull(), utils.TimeZero(), utils.SQLNullTimeNull(),
				)
				do := Orders(m)
				order := entity.NewOrder("test", "0")
				return do.Insert(order)
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
