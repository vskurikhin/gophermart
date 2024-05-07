/*
 * This file was last modified at 2024-05-07 17:40 by Victor N. Skurikhin.
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

// Handle balance
//
//	@Summary		баланс
//	@Description	получение текущего баланса пользователя
//	@Tags			Account
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	model.Balance	"успешная обработка запроса"
//	@Failure		401	{object}	model.JSONError	"пользователь не аутентифицирован"
//	@Failure		500	{string}	string			"Internal Server Error"
//	@Router			/user/balance/withdraw [get]
//
//goland:noinspection GoUnhandledErrorResult
func (r *balance) Handle(response http.ResponseWriter, request *http.Request) {

	ctx := request.Context()
	login, err := auth.Login(ctx)

	if err != nil || login == nil {

		r.log.Debug(balanceMsg, utils.LogCtxRecoverFields(ctx, err)...)
		render.Status(request, http.StatusBadRequest)
		render.Render(response, request, model.Error(handlers.ErrBadRequest))

		return
	}
	result := newService(ctx, storage.NewPgsStorage()).Balance(*login)
	render.Status(request, result.Status())

	switch value := result.(type) {
	case *handlers.ResultError:
		render.Render(response, request, model.Error(value.Error()))
	case *handlers.ResultAny:

		if balance, ok := value.Result().(*model.Balance); ok {
			response.Header().Set("Content-Type", "application/json")
			if err := balance.MarshalToWriter(response); err == nil {
				return
			} else {
				r.log.Debug(balanceMsg, utils.LogCtxReasonErrFields(ctx, err.Error(), handlers.ErrInternalError)...)
			}
		}
	}
	r.log.Debug(balanceMsg, utils.InternalErrorZapField(ctx, request, result)...)
}
