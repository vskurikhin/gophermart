/*
 * This file was last modified at 2024-04-16 10:44 by Victor N. Skurikhin.
 * register.go
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

const regMsg = "register"

type register struct {
	log *zap.Logger
}

func newRegister() *register {
	return &register{log: logger.Get()}
}

func (r *register) Handle(response http.ResponseWriter, request *http.Request) {

	hr := newHandleResult(response, request)
	hr.handleUser(logMsg, func(ctx context.Context, userRegister *model.User) handlers.Result {
		return NewUserService(ctx).Register(userRegister)
	})
}
