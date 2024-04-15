/*
 * This file was last modified at 2024-04-15 10:39 by Victor N. Skurikhin.
 * holder.go
 * $Id$
 */

package accounts

import (
	"github.com/vskurikhin/gophermart/internal/handlers"
	"sync"
)

type Holder interface {
	HandlerBalance() handlers.Handler
	HandlerBalanceWithdraw() handlers.Handler
	HandlerWithdrawals() handlers.Handler
}

type holder struct {
	balance         handlers.Handler
	balanceWithdraw handlers.Handler
	withdrawals     handlers.Handler
}

var once = new(sync.Once)
var instance Holder

func GetInstance() Holder {

	once.Do(func() {
		h := new(holder)
		h.balance = newBalance()
		h.balanceWithdraw = newBalanceWithdraw()
		h.withdrawals = newWithdraws()
		instance = h
	})
	return instance
}

func (h *holder) HandlerBalance() handlers.Handler {
	return h.balance
}

func (h *holder) HandlerBalanceWithdraw() handlers.Handler {
	return h.balanceWithdraw
}

func (h *holder) HandlerWithdrawals() handlers.Handler {
	return h.withdrawals
}
