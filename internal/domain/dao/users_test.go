/*
 * This file was last modified at 2024-04-22 10:44 by Victor N. Skurikhin.
 * users_test.go
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
func TestDaoUsers(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var tests = []struct {
		name string
		get  func() (interface{}, error)
	}{
		{
			name: "Test #1 Users.GetUser",
			get: func() (interface{}, error) {
				columns := []string{"login", "password", "created_at", "update_at"}
				m := MockStorageGetByString(
					ctrl, columns, "test", utils.SQLNullStringNull(), utils.TimeZero(), utils.SQLNullTimeNull(),
				)
				du := Users(m)
				return du.GetUser("test")
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
