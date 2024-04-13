/*
 * This file was last modified at 2024-04-13 17:14 by Victor N. Skurikhin.
 * balances.go
 * $Id$
 */

package dao

import (
	"github.com/vskurikhin/gophermart/internal/domain/entity"
	"github.com/vskurikhin/gophermart/internal/storage"
)

type balances struct {
	storage storage.Storage
}

func Balances(storage storage.Storage) *balances {
	return &balances{storage: storage}
}

func (b *balances) GetAllBalances() ([]*entity.Balance, error) {
	return entity.FuncGetAllBalances()(b.storage)
}

func (b *balances) GetBalance(login string) (*entity.Balance, error) {
	return entity.FuncGetBalance()(b.storage, login)
}

func (b *balances) Save(balance *entity.Balance) (*entity.Balance, error) {
	return balance.Save(b.storage)
}
