/*
 * This file was last modified at 2024-04-19 14:39 by Victor N. Skurikhin.
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
