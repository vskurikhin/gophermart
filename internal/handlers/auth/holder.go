/*
 * This file was last modified at 2024-04-15 14:48 by Victor N. Skurikhin.
 * holder.go
 * $Id$
 */

package auth

import (
	"github.com/vskurikhin/gophermart/internal/handlers"
	"github.com/vskurikhin/gophermart/internal/logger"
	"github.com/vskurikhin/gophermart/internal/storage"
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
		h.login = newLogin(logger.Get(), storage.NewPgsStorage())
		h.register = newRegister(logger.Get(), storage.NewPgsStorage())
		instance = h
	})
	return instance
}
