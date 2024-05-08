/*
 * This file was last modified at 2024-05-07 17:40 by Victor N. Skurikhin.
 * balance_withdraw.go
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

const balanceWithdrawMsg = "balanceWithdraw"

type balanceWithdraw struct {
	log *zap.Logger
}

func newBalanceWithdraw() *balanceWithdraw {
	return &balanceWithdraw{log: logger.Get()}
}

// Handle balanceWithdraw
//
//	@Summary		списание
//	@Description	запрос на списание средств
//	@Tags			Account
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		model.Withdraw	true	"Формат запроса JSON (body)"
//	@Success		200		{object}	model.Withdraw	"успешная обработка запроса"
//	@Failure		401		{object}	model.JSONError	"пользователь не аутентифицирован"
//	@Failure		402		{object}	model.JSONError	"на счету недостаточно средств"
//	@Failure		422		{object}	model.JSONError	"неверный номер заказа"
//	@Failure		500		{string}	string			"Internal Server Error"
//	@Router			/user/balance/withdraw [post]
//
//goland:noinspection GoUnhandledErrorResult
func (r *balanceWithdraw) Handle(response http.ResponseWriter, request *http.Request) {

	ctx := request.Context()
	login, err := auth.Login(ctx)

	if err != nil || login == nil {

		r.log.Debug(balanceWithdrawMsg, utils.LogCtxRecoverFields(ctx, err)...)
		render.Status(request, http.StatusBadRequest)
		render.Render(response, request, model.Error(handlers.ErrBadRequest))

		return
	}
	withdraw, err := model.UnmarshalWithdrawFromReader(request.Body)

	if err != nil || withdraw == nil {

		r.log.Debug(balanceWithdrawMsg, utils.LogCtxRecoverFields(ctx, err)...)
		render.Status(request, http.StatusBadRequest)
		render.Render(response, request, model.Error(handlers.ErrBadRequest))

		return
	}
	result := newService(ctx, storage.NewPgsStorage()).Withdraw(*login, withdraw)
	render.Status(request, result.Status())

	switch value := result.(type) {
	case *handlers.ResultError:
		render.Render(response, request, model.Error(value.Error()))
	case *handlers.ResultAny:

		if withdraw, ok := value.Result().(*model.Withdraw); ok {
			response.Header().Set("Content-Type", "application/json")
			if err := withdraw.MarshalToWriter(response); err == nil {
				return
			} else {
				r.log.Debug(balanceWithdrawMsg, utils.LogCtxReasonErrFields(ctx, err.Error(), handlers.ErrInternalError)...)
			}
		}
	}
	r.log.Debug(balanceWithdrawMsg, utils.InternalErrorZapField(ctx, request, result)...)
}
