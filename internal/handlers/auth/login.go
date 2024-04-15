/*
 * This file was last modified at 2024-04-15 16:54 by Victor N. Skurikhin.
 * login.go
 * $Id$
 */

package auth

import (
	"github.com/go-chi/render"
	"github.com/vskurikhin/gophermart/internal/handlers"
	"github.com/vskurikhin/gophermart/internal/model"
	"github.com/vskurikhin/gophermart/internal/storage"
	"go.uber.org/zap"
	"net/http"
)

type login struct {
	log   *zap.Logger
	store *storage.PgsStorage
}

func newLogin(log *zap.Logger, store *storage.PgsStorage) *login {
	return &login{log: log, store: store}
}

func (r *login) Handle(response http.ResponseWriter, request *http.Request) {

	userRegister, err := model.UnmarshalFromReader(request.Body)

	if err != nil {
		render.Status(request, http.StatusBadRequest)
		//goland:noinspection GoUnhandledErrorResult
		render.Render(response, request, model.Error(handlers.ErrBadRequest))
	}
	if err := userRegister.MarshalToWriter(response); err != nil {
		panic(err)
	}
}
