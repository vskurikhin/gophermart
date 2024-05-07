/*
 * This file was last modified at 2024-05-07 17:40 by Victor N. Skurikhin.
 * register.go
 * $Id$
 */

package auth

import (
	"context"
	"github.com/vskurikhin/gophermart/internal/handlers"
	"github.com/vskurikhin/gophermart/internal/logger"
	"github.com/vskurikhin/gophermart/internal/model"
	"github.com/vskurikhin/gophermart/internal/storage"
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

// Handle register
//
//	@Summary		регистрация
//	@Description	регистрация пользователя
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request			body		model.User		true	"Формат запроса JSON (body)"
//	@Success		200				{object}	model.User		"пользователь успешно зарегистрирован и аутентифицирован"
//	@Failure		400				{object}	model.JSONError	"неверный формат запроса"
//	@Failure		409				{object}	model.JSONError	"логин уже занят"
//	@Failure		500				{string}	string			"Internal Server Error"
//	@Router			/user/register 																								[post]
func (r *register) Handle(response http.ResponseWriter, request *http.Request) {

	hr := newHandleResult(response, request)
	hr.handleUser(logMsg, func(ctx context.Context, userRegister *model.User) handlers.Result {
		return NewUserService(ctx, storage.NewPgsStorage()).Register(userRegister)
	})
}
