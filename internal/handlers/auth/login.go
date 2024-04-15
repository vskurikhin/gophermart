/*
 * This file was last modified at 2024-04-16 10:44 by Victor N. Skurikhin.
 * login.go
 * $Id$
 */

package auth

import (
	"context"
	"github.com/vskurikhin/gophermart/internal/handlers"
	"github.com/vskurikhin/gophermart/internal/logger"
	"github.com/vskurikhin/gophermart/internal/model"
	"go.uber.org/zap"
	"net/http"
)

const logMsg = "login"

type login struct {
	log *zap.Logger
}

func newLogin() *login {
	return &login{log: logger.Get()}
}

func (l *login) Handle(response http.ResponseWriter, request *http.Request) {

	hr := newHandleResult(response, request)
	hr.handleUser(logMsg, func(ctx context.Context, userRegister *model.User) handlers.Result {
		return NewUserService(ctx).Login(userRegister)
	})
}
