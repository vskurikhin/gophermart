/*
 * This file was last modified at 2024-04-15 14:48 by Victor N. Skurikhin.
 * holder.go
 * $Id$
 */

package accounts

import (
	"github.com/vskurikhin/gophermart/internal/handlers"
	"sync"
)

type holder struct {
	balance         handlers.Handler
	balanceWithdraw handlers.Handler
	withdrawals     handlers.Handler
}

var once = new(sync.Once)
var instance *holder

func getInstance() *holder {

	once.Do(func() {
		h := new(holder)
		h.balance = newBalance()
		h.balanceWithdraw = newBalanceWithdraw()
		h.withdrawals = newWithdraws()
		instance = h
	})
	return instance
}
