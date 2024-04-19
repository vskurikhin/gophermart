/*
 * This file was last modified at 2024-04-19 23:20 by Victor N. Skurikhin.
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

func (r *balanceWithdraw) Handle(response http.ResponseWriter, request *http.Request) {

	ctx := request.Context()
	login, err := auth.Login(ctx)

	if err != nil || login == nil {

		r.log.Debug(balanceWithdrawMsg, utils.LogCtxRecoverFields(ctx, err)...)
		render.Status(request, http.StatusBadRequest)
		//goland:noinspection GoUnhandledErrorResult
		render.Render(response, request, model.Error(handlers.ErrBadRequest))

		return
	}
	withdraw, err := model.UnmarshalWithdrawFromReader(request.Body)

	if err != nil || withdraw == nil {

		r.log.Debug(balanceWithdrawMsg, utils.LogCtxRecoverFields(ctx, err)...)
		render.Status(request, http.StatusBadRequest)
		//goland:noinspection GoUnhandledErrorResult
		render.Render(response, request, model.Error(handlers.ErrBadRequest))

		return
	}
	result := newService(ctx, storage.NewPgsStorage()).Withdraw(*login, withdraw)

	switch value := result.(type) {
	case *handlers.ResultError:

		render.Status(request, value.Status())
		//goland:noinspection GoUnhandledErrorResult
		render.Render(response, request, model.Error(value.Error()))

	case *handlers.ResultAny:

		if withdraw, ok := value.Result().(*model.Withdraw); ok {
			response.Header().Set("Content-Type", "application/json")
			if err := withdraw.MarshalToWriter(response); err == nil {
				render.Status(request, value.Status())
				return
			} else {
				r.log.Debug(balanceWithdrawMsg, utils.LogCtxReasonErrFields(ctx, err.Error(), handlers.ErrInternalError)...)
			}
		}
	}
	r.log.Debug(balanceWithdrawMsg, utils.InternalErrorZapField(ctx, request, result)...)

}
