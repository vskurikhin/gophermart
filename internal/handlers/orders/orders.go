/*
 * This file was last modified at 2024-04-20 18:04 by Victor N. Skurikhin.
 * orders.go
 * $Id$
 */

package orders

import (
	"github.com/go-chi/render"
	"github.com/vskurikhin/gophermart/internal/clients/accrual"
	"github.com/vskurikhin/gophermart/internal/handlers"
	"github.com/vskurikhin/gophermart/internal/handlers/auth"
	"github.com/vskurikhin/gophermart/internal/logger"
	"github.com/vskurikhin/gophermart/internal/model"
	"github.com/vskurikhin/gophermart/internal/storage"
	"github.com/vskurikhin/gophermart/internal/utils"
	"go.uber.org/zap"
	"net/http"
)

const ordersMsg = "orders"

type orders struct {
	log *zap.Logger
}

func newOrders() *orders {
	return &orders{log: logger.Get()}
}

//goland:noinspection GoUnhandledErrorResult
func (o *orders) Handle(response http.ResponseWriter, request *http.Request) {

	ctx := request.Context()
	login, err := auth.Login(ctx)

	if err != nil || login == nil {

		o.log.Debug(ordersMsg, utils.LogCtxRecoverFields(ctx, err)...)
		render.Status(request, http.StatusBadRequest)
		render.Render(response, request, model.Error(handlers.ErrBadRequest))

		return
	}
	result := newService(ctx, storage.NewPgsStorage(), accrual.GetWorkers()).
		Orders(*login)
	render.Status(request, result.Status())

	switch value := result.(type) {
	case *handlers.ResultError:
		render.Render(response, request, model.Error(value.Error()))

	case *handlers.ResultAny:

		if orders, ok := value.Result().(model.Orders); ok {
			response.Header().Set("Content-Type", "application/json")
			if err := orders.MarshalToWriter(response); err == nil {
				return
			} else {
				o.log.Debug(ordersMsg, utils.LogCtxReasonErrFields(ctx, err.Error(), handlers.ErrInternalError)...)
			}
		}
	}
	o.log.Debug(ordersMsg, utils.InternalErrorZapField(ctx, request, result)...)
}
