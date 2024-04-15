/*
 * This file was last modified at 2024-04-15 14:48 by Victor N. Skurikhin.
 * holder.go
 * $Id$
 */

package orders

import (
	"github.com/vskurikhin/gophermart/internal/handlers"
	"sync"
)

type holder struct {
	number handlers.Handler
	orders handlers.Handler
}

var once = new(sync.Once)
var instance *holder

func getInstance() *holder {

	once.Do(func() {
		h := new(holder)
		h.number = newNumber()
		h.orders = newOrders()
		instance = h
	})
	return instance
}
