/*
 * This file was last modified at 2024-05-07 17:21 by Victor N. Skurikhin.
 * login.go
 * $Id$
 */

package auth

import (
	"context"
	"fmt"
	"github.com/go-chi/jwtauth"
	"github.com/vskurikhin/gophermart/internal/handlers"
	"github.com/vskurikhin/gophermart/internal/logger"
	"github.com/vskurikhin/gophermart/internal/model"
	"github.com/vskurikhin/gophermart/internal/storage"
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

// Handle login
//
//	@Summary		аутентификация
//	@Description	аутентификация пользователя
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Security		none
//	@Param			request			body		model.User		true	"Формат запроса JSON (body)"
//	@Success		200				{object}	model.User		"пользователь успешно аутентифицирован"
//	@Failure		400				{object}	model.JSONError	"неверный формат запроса"
//	@Failure		401				{object}	model.JSONError	"неверная пара логин/пароль"
//	@Failure		500				{string}	string			"Internal Server Error"
//	@Router			/user/login 	[post]
func (l *login) Handle(response http.ResponseWriter, request *http.Request) {

	hr := newHandleResult(response, request)
	hr.handleUser(logMsg, func(ctx context.Context, userRegister *model.User) handlers.Result {
		return NewUserService(ctx, storage.NewPgsStorage()).Login(userRegister)
	})
}

func Login(ctx context.Context) (*string, error) {

	_, m, err := jwtauth.FromContext(ctx)

	if err != nil {
		return nil, err
	}
	if login, ok := m["username"].(string); ok {
		return &login, nil
	}
	return nil, fmt.Errorf("username not found")
}
