/*
 * This file was last modified at 2024-05-13 15:47 by Victor N. Skurikhin.
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

// Handle withdrawals
//
//	@Summary		информация о выводе средств
//	@Description	получение информации о выводе средств
//	@Tags			Account
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200				{array}		model.Withdraw	"успешная обработка запроса"
//	@Success		204				{object}	model.Empty		"нет ни одного списания"
//	@Failure		401				{object}	model.JSONError	"пользователь не авторизован"
//	@Failure		500				{string}	string			"Internal Server Error"
//	@Router			/user/withdrawals 	[get]
//
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
