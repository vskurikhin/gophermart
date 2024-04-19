/*
 * This file was last modified at 2024-04-19 23:19 by Victor N. Skurikhin.
 * balance.go
 * $Id$
 */

package accounts

import (
	"github.com/go-chi/render"
	"github.com/vskurikhin/gophermart/internal/handlers"
	"github.com/vskurikhin/gophermart/internal/handlers/auth"
	"github.com/vskurikhin/gophermart/internal/logger"
	"github.com/vskurikhin/gophermart/internal/model"
	"github.com/vskurikhin/gophermart/internal/storage"
	"github.com/vskurikhin/gophermart/internal/utils"
	"go.uber.org/zap"
	"net/http"
)

const balanceMsg = "balance"

type balance struct {
	log *zap.Logger
}

func newBalance() *balance {
	return &balance{log: logger.Get()}
}

func (r *balance) Handle(response http.ResponseWriter, request *http.Request) {

	ctx := request.Context()
	login, err := auth.Login(ctx)

	if err != nil || login == nil {

		r.log.Debug(balanceMsg, utils.LogCtxRecoverFields(ctx, err)...)
		render.Status(request, http.StatusBadRequest)
		//goland:noinspection GoUnhandledErrorResult
		render.Render(response, request, model.Error(handlers.ErrBadRequest))

		return
	}
	result := newService(ctx, storage.NewPgsStorage()).Balance(*login)

	switch value := result.(type) {
	case *handlers.ResultError:

		render.Status(request, value.Status())
		//goland:noinspection GoUnhandledErrorResult
		render.Render(response, request, model.Error(value.Error()))

	case *handlers.ResultAny:

		if balance, ok := value.Result().(*model.Balance); ok {
			response.Header().Set("Content-Type", "application/json")
			if err := balance.MarshalToWriter(response); err == nil {
				render.Status(request, value.Status())
				return
			} else {
				r.log.Debug(balanceMsg, utils.LogCtxReasonErrFields(ctx, err.Error(), handlers.ErrInternalError)...)
			}
		}
	}
	r.log.Debug(balanceMsg, utils.InternalErrorZapField(ctx, request, result)...)
}
