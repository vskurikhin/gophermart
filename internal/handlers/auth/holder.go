/*
 * This file was last modified at 2024-04-16 10:37 by Victor N. Skurikhin.
 * holder.go
 * $Id$
 */

package auth

import (
	"github.com/vskurikhin/gophermart/internal/handlers"
	"net/http"
	"sync"
)

func UserLoginHandlerFunc() http.HandlerFunc {
	return http.HandlerFunc(getInstance().login.Handle)
}

func UserRegisterHandlerFunc() http.HandlerFunc {
	return http.HandlerFunc(getInstance().register.Handle)
}

type holder struct {
	login    handlers.Handler
	register handlers.Handler
}

var once = new(sync.Once)
var instance *holder

func getInstance() *holder {

	once.Do(func() {
		h := new(holder)
		h.login = newLogin()
		h.register = newRegister()
		instance = h
	})
	return instance
}
