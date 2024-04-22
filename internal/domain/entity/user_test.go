/*
 * This file was last modified at 2024-04-21 15:11 by Victor N. Skurikhin.
 * user_test.go
 * $Id$
 */

package entity

import (
	"github.com/stretchr/testify/assert"
	"github.com/vskurikhin/gophermart/internal/storage"
	"github.com/vskurikhin/gophermart/internal/utils"
	"testing"
)

var userColumns = []string{"login", "password", "created_at", "update_at"}

//goland:noinspection GoImportUsedAsName,GoUnhandledErrorResult
func TestUserPositive(t *testing.T) {

	var tests = []struct {
		name string
		get  func() (interface{}, error)
	}{
		{
			name: "Test positive #1 User",
			get: func() (interface{}, error) {

				u := NewUser("test", nil)
				assert.NotNil(t, u)
				assert.Equal(t, "test", u.Login())
				assert.Nil(t, u.Password())

				a := make(storage.TxArgs, 0)
				assert.True(t, len(a) == 0)

				a = u.AppendInsertTo(a)
				assert.True(t, len(a) == 1)

				rows := utils.PgxRowsNext(
					userColumns, "test", utils.SqlNullStringZero(), utils.TimeZero(), utils.SqlNullTimeZero(),
				)
				return extractUser(rows)
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
func TestUserNegative(t *testing.T) {

	var tests = []struct {
		name string
		get  func() (interface{}, error)
	}{
		{
			name: "Test negative #1 User",
			get: func() (interface{}, error) {
				rows := &utils.TestRows{}
				return extractUser(rows)
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
