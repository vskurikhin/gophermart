/*
 * This file was last modified at 2024-04-21 12:38 by Victor N. Skurikhin.
 * mock_test.go
 * $Id$
 */

package dao

import (
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/vskurikhin/gophermart/internal/utils"
	"go.uber.org/mock/gomock"
)

func MockStorageGetByString(ctrl *gomock.Controller, columns []string, values ...interface{}) *MockStorage {

	m := NewMockStorage(ctrl)
	m.EXPECT().
		GetByString(gomock.Any(), gomock.Any()).
		DoAndReturn(func(s1 string, s2 string) (pgx.Row, error) {
			return utils.PgxRowsNext(columns, values...), nil
		}).AnyTimes()
	return m
}

func MockStorageGetAllForString(ctrl *gomock.Controller, columns []string, values ...interface{}) *MockStorage {

	m := NewMockStorage(ctrl)
	m.EXPECT().
		GetAllForString(gomock.Any(), gomock.Any()).
		DoAndReturn(func(s1 string, s2 string) (interface{}, error) {
			return nil, fmt.Errorf("unimpl")
		}).AnyTimes()
	return m
}

func MockStorageSaveOrder(ctrl *gomock.Controller, columns []string, values ...interface{}) *MockStorage {

	m := NewMockStorage(ctrl)
	vs := []any{gomock.Any(), gomock.Any()}
	m.EXPECT().
		Save(gomock.Any(), vs...).
		DoAndReturn(func(s1 string, s2 ...any) (pgx.Row, error) {
			return utils.PgxRowsNext(columns, values...), nil
		}).AnyTimes()
	return m
}
