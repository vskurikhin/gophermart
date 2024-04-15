/*
 * This file was last modified at 2024-04-15 10:39 by Victor N. Skurikhin.
 * holder.go
 * $Id$
 */

package auth

import (
	"github.com/vskurikhin/gophermart/internal/handlers"
	"sync"
)

type Holder interface {
	HandlerLogin() handlers.Handler
	HandlerRegister() handlers.Handler
}

type holder struct {
	login    handlers.Handler
	register handlers.Handler
}

var once = new(sync.Once)
var instance Holder

func GetInstance() Holder {

	once.Do(func() {
		h := new(holder)
		h.login = newLogin()
		h.register = newRegister()
		instance = h
	})
	return instance
}

func (h *holder) HandlerLogin() handlers.Handler {
	return h.login
}

func (h *holder) HandlerRegister() handlers.Handler {
	return h.register
}
