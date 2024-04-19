/*
 * This file was last modified at 2024-04-20 00:53 by Victor N. Skurikhin.
 * withdrawals.go
 * $Id$
 */

package dao

import (
	"github.com/vskurikhin/gophermart/internal/domain/entity"
	"github.com/vskurikhin/gophermart/internal/storage"
	"math/big"
)

type withdrawals struct {
	storage storage.Storage
}

func Withdrawals(storage storage.Storage) *withdrawals {
	return &withdrawals{storage: storage}
}

func (w *withdrawals) GetAllWithdraws() ([]*entity.Withdraw, error) {
	return entity.FuncGetAllWithdraw()(w.storage)
}

func (w *withdrawals) GetAllWithdrawalsByLogin(login string) ([]*entity.Withdraw, error) {
	return entity.FuncGetAllWithdrawalsByLogin()(w.storage, login)
}

func (w *withdrawals) GetWithdraw(login, number string) (*entity.Withdraw, error) {
	return entity.FuncGetWithdraw()(w.storage, login, number)
}

func (w *withdrawals) GetWithdrawSum(login string) (*big.Float, error) {
	return entity.FuncGetWithdrawSum()(w.storage, login)
}

func (w *withdrawals) Save(withdraw *entity.Withdraw) (*entity.Withdraw, error) {
	return withdraw.Save(w.storage)
}
