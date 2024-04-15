/*
 * This file was last modified at 2024-04-15 23:37 by Victor N. Skurikhin.
 * db_test.go
 * $Id$
 */

package dao

import (
	"github.com/stretchr/testify/assert"
	"github.com/vskurikhin/gophermart/internal/domain/entity"
	"github.com/vskurikhin/gophermart/internal/storage"
	"github.com/vskurikhin/gophermart/internal/utils"
	"math/big"
	"testing"
)

//goland:noinspection GoImportUsedAsName,GoUnhandledErrorResult
func TestDaoSaveGet(t *testing.T) {

	var store storage.Storage
	store = storage.NewPgsStorage()

	var tests = []struct {
		name string
		repo interface{}
		save func() (interface{}, error)
		get  func() (interface{}, error)
	}{
		{
			name: "Test #1 User",
			save: func() (interface{}, error) {
				store = store.WithContext(utils.NewIDContext())
				repo0 := Users(store)
				test0 := entity.NewUser("test1", nil)
				repo0.Insert(test0)
				return repo0.Save(test0)
			},
			get: func() (interface{}, error) {
				store = store.WithContext(utils.NewIDContext())
				repo0 := Users(store)
				return repo0.GetUser("test1")
			},
		},
		{
			name: "Test #2 Balance",
			save: func() (interface{}, error) {
				store = store.WithContext(utils.NewIDContext())
				repo0 := Balances(store)
				test0 := entity.NewBalance("test1", *big.NewFloat(0.1))
				return repo0.Save(test0)
			},
			get: func() (interface{}, error) {
				store = store.WithContext(utils.NewIDContext())
				repo0 := Balances(store)
				return repo0.GetBalance("test1")
			},
		},
		{
			name: "Test #3 Order",
			save: func() (interface{}, error) {
				store = store.WithContext(utils.NewIDContext())
				repo0 := Orders(store)
				test0 := entity.NewOrder("test1", "test1", 1)
				return repo0.Save(test0)
			},
			get: func() (interface{}, error) {
				store = store.WithContext(utils.NewIDContext())
				repo0 := Orders(store)
				return repo0.GetOrder("test1", "test1")
			},
		},
		{
			name: "Test #4 Withdraws",
			save: func() (interface{}, error) {
				store = store.WithContext(utils.NewIDContext())
				repo0 := Withdraws(store)
				test0 := entity.NewWithdraw("test1", "test1", *new(big.Float), 1)
				return repo0.Save(test0)
			},
			get: func() (interface{}, error) {
				store = store.WithContext(utils.NewIDContext())
				repo0 := Withdraws(store)
				return repo0.GetWithdraw("test1", "test1")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test1, err := test.save()
			assert.Nil(t, err)
			test2, err := test.get()
			assert.Nil(t, err)
			assert.Equal(t, test1, test2)
		})
	}
}

//goland:noinspection GoImportUsedAsName,GoUnhandledErrorResult
func TestStatus(t *testing.T) {

	var tests = []struct {
		name string
	}{
		{"Test Status #1"},
	}

	var store storage.Storage
	store = storage.NewPgsStorage()
	store = store.WithContext(utils.NewIDContext())
	repo0 := Statuses(store)
	test0 := &entity.Status{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := repo0.Save(test0)
			assert.Nil(t, err)
			test1, err := repo0.GetStatus(1)
			assert.Nil(t, err)
			test_, err := repo0.GetStatus(1)
			assert.Nil(t, err)
			assert.Equal(t, test1, test_)
			repo1 := Statuses(store)
			assert.Equal(t, repo0, repo1)
			repo2 := Statuses(store.WithContext(utils.NewIDContext()))
			assert.NotEqual(t, repo0, repo2)
		})
	}
}
