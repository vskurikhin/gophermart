/*
 * This file was last modified at 2024-04-20 17:09 by Victor N. Skurikhin.
 * balance_order.go
 * $Id$
 */

package transaction

import (
	"github.com/vskurikhin/gophermart/internal/domain/entity"
	"github.com/vskurikhin/gophermart/internal/storage"
)

type balanceOrder struct {
	store storage.Storage
}

func BalanceOrder(store storage.Storage) *balanceOrder {
	return &balanceOrder{store: store}
}

func (u *balanceOrder) TransactionAccrual(balance *entity.Balance, order *entity.Order, status string) error {

	args := make(storage.TxArgs, 0)
	args = order.AppendAccrualTo(args)
	args = order.AppendSetStatusTo(args, status)

	if order.Accrual() != nil {
		args = balance.AppendAccrualTo(args, order.Accrual())
	}
	return u.store.Transaction(args)
}
