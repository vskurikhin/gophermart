/*
 * This file was last modified at 2024-04-20 18:07 by Victor N. Skurikhin.
 * withdrawals.go
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

const withdrawalsMsg = "balanceWithdraw"

type withdrawals struct {
	log *zap.Logger
}

func newWithdraws() *withdrawals {
	return &withdrawals{log: logger.Get()}
}

//goland:noinspection GoUnhandledErrorResult
func (r *withdrawals) Handle(response http.ResponseWriter, request *http.Request) {

	ctx := request.Context()
	login, err := auth.Login(ctx)

	if err != nil || login == nil {

		r.log.Debug(withdrawalsMsg, utils.LogCtxRecoverFields(ctx, err)...)
		render.Status(request, http.StatusBadRequest)
		render.Render(response, request, model.Error(handlers.ErrBadRequest))

		return
	}
	result := newService(ctx, storage.NewPgsStorage()).Withdrawals(*login)
	render.Status(request, result.Status())

	switch value := result.(type) {
	case *handlers.ResultError:
		render.Render(response, request, model.Error(value.Error()))
	case *handlers.ResultAny:

		if withdrawals, ok := value.Result().(model.Withdrawals); ok {
			response.Header().Set("Content-Type", "application/json")
			if err := withdrawals.MarshalToWriter(response); err == nil {
				return
			} else {
				r.log.Debug(withdrawalsMsg, utils.LogCtxReasonErrFields(ctx, err.Error(), handlers.ErrInternalError)...)
			}
		}
	case *handlers.ResultString:
		response.Header().Set("Content-Type", "application/json")
		render.Render(response, request, model.NewEmptyList())
		return
	}
	r.log.Debug(withdrawalsMsg, utils.InternalErrorZapField(ctx, request, result)...)
}
