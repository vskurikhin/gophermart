/*
 * This file was last modified at 2024-04-20 17:50 by Victor N. Skurikhin.
 * balances.go
 * $Id$
 */

package dao

import (
	"github.com/vskurikhin/gophermart/internal/domain/entity"
	"github.com/vskurikhin/gophermart/internal/domain/view"
	"github.com/vskurikhin/gophermart/internal/storage"
)

type balances struct {
	storage storage.Storage
}

func Balances(storage storage.Storage) *balances {
	return &balances{storage: storage}
}

func (b *balances) GetBalance(login string) (*entity.Balance, error) {
	return entity.FuncGetBalance()(b.storage, login)
}

func (b *balances) GetBalanceWithdraw(login string) (*view.Balance, error) {
	return view.FuncGetBalanceWithdraw()(b.storage, login)
}
