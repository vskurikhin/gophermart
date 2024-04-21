/*
 * This file was last modified at 2024-04-19 21:04 by Victor N. Skurikhin.
 * user_balance.go
 * $Id$
 */

package transaction

import (
	"github.com/vskurikhin/gophermart/internal/domain/entity"
	"github.com/vskurikhin/gophermart/internal/storage"
)

type userBalance struct {
	store storage.Storage
}

func UserBalance(store storage.Storage) *userBalance {
	return &userBalance{store: store}
}

func (u *userBalance) TransactionInsert(user *entity.User, balance *entity.Balance) error {

	args := make(storage.TxArgs, 0)
	args = user.AppendInsertTo(args)
	args = balance.AppendInsertTo(args)

	return u.store.Transaction(args)
}
