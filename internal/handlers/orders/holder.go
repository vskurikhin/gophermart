/*
 * This file was last modified at 2024-04-15 10:39 by Victor N. Skurikhin.
 * holder.go
 * $Id$
 */

package orders

import (
	"github.com/vskurikhin/gophermart/internal/handlers"
	"sync"
)

type Holder interface {
	HandlerNumber() handlers.Handler
	HandlerOrders() handlers.Handler
}

type holder struct {
	number handlers.Handler
	orders handlers.Handler
}

var once = new(sync.Once)
var instance Holder

func GetInstance() Holder {

	once.Do(func() {
		h := new(holder)
		h.number = newNumber()
		h.orders = newOrders()
		instance = h
	})
	return instance
}

func (h *holder) HandlerNumber() handlers.Handler {
	return h.number
}

func (h *holder) HandlerOrders() handlers.Handler {
	return h.orders
}
