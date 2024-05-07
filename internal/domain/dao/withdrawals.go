/*
 * This file was last modified at 2024-04-25 21:59 by Victor N. Skurikhin.
 * withdrawals.go
 * $Id$
 */

package dao

import (
	"github.com/vskurikhin/gophermart/internal/domain/entity"
	"github.com/vskurikhin/gophermart/internal/storage"
)

type withdrawals struct {
	storage storage.Storage
}

func Withdrawals(storage storage.Storage) *withdrawals {
	return &withdrawals{storage: storage}
}

func (w *withdrawals) GetAllWithdrawalsByLogin(login string) ([]*entity.Withdraw, error) {
	return entity.GetAllWithdrawalsByLogin(w.storage, login)
}
