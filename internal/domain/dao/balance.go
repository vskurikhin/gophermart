/*
 * This file was last modified at 2024-04-12 16:18 by Victor N. Skurikhin.
 * balance.go
 * $Id$
 */

package dao

import (
	"github.com/vskurikhin/gophermart/internal/domain/entity"
	"github.com/vskurikhin/gophermart/internal/storage"
)

type balance struct {
	storage storage.Storage
}

func Balances(storage storage.Storage) *balance {
	return &balance{storage: storage}
}

func (w *balance) GetAllBalances() ([]*entity.Balance, error) {
	return entity.FuncGetAllBalances()(w.storage)
}

func (w *balance) GetBalance(login string) (*entity.Balance, error) {
	return entity.FuncGetBalance()(w.storage, login)
}

func (w *balance) Save(balance *entity.Balance) (*entity.Balance, error) {
	return balance.Save(w.storage)
}
