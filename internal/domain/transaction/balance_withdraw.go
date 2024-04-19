/*
 * This file was last modified at 2024-04-19 22:38 by Victor N. Skurikhin.
 * balance_withdraw.go
 * $Id$
 */

package transaction

import (
	"github.com/vskurikhin/gophermart/internal/domain/entity"
	"github.com/vskurikhin/gophermart/internal/storage"
)

type balanceWithdraw struct {
	store storage.Storage
}

func BalanceWithdraw(store storage.Storage) *balanceWithdraw {
	return &balanceWithdraw{store: store}
}

func (u *balanceWithdraw) TransactionWithdraw(balance *entity.Balance, withdraw *entity.Withdraw) error {

	args := make(storage.TxArgs, 0)
	args = withdraw.AppendInsertTo(args)
	args = balance.AppendWithdrawTo(args, withdraw.Sum())

	return u.store.Transaction(args)
}
