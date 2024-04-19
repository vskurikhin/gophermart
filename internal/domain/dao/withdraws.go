/*
 * This file was last modified at 2024-04-19 15:58 by Victor N. Skurikhin.
 * withdraws.go
 * $Id$
 */

package dao

import (
	"github.com/vskurikhin/gophermart/internal/domain/entity"
	"github.com/vskurikhin/gophermart/internal/storage"
	"math/big"
)

type withdraws struct {
	storage storage.Storage
}

func Withdraws(storage storage.Storage) *withdraws {
	return &withdraws{storage: storage}
}

func (w *withdraws) GetAllWithdraws() ([]*entity.Withdraw, error) {
	return entity.FuncGetAllWithdraw()(w.storage)
}

func (w *withdraws) GetWithdraw(login, number string) (*entity.Withdraw, error) {
	return entity.FuncGetWithdraw()(w.storage, login, number)
}

func (w *withdraws) GetWithdrawSum(login string) (*big.Float, error) {
	return entity.FuncGetWithdrawSum()(w.storage, login)
}

func (w *withdraws) Save(withdraw *entity.Withdraw) (*entity.Withdraw, error) {
	return withdraw.Save(w.storage)
}
