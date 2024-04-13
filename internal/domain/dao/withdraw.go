/*
 * This file was last modified at 2024-04-12 16:01 by Victor N. Skurikhin.
 * withdraw.go
 * $Id$
 */

package dao

import (
	"github.com/vskurikhin/gophermart/internal/domain/entity"
	"github.com/vskurikhin/gophermart/internal/storage"
)

type withdraw struct {
	storage storage.Storage
}

func Withdraw(storage storage.Storage) *withdraw {
	return &withdraw{storage: storage}
}

func (w *withdraw) GetAllWithdraw() ([]*entity.Withdraw, error) {
	return entity.FuncGetAllWithdraw()(w.storage)
}

func (w *withdraw) GetWithdraw(login, number string) (*entity.Withdraw, error) {
	return entity.FuncGetWithdraw()(w.storage, login, number)
}

func (w *withdraw) Save(withdraw *entity.Withdraw) (*entity.Withdraw, error) {
	return withdraw.Save(w.storage)
}
