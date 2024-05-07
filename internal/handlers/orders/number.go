/*
 * This file was last modified at 2024-05-07 17:19 by Victor N. Skurikhin.
 * number.go
 * $Id$
 */

package orders

import (
	"github.com/go-chi/render"
	"github.com/vskurikhin/gophermart/internal/clients/accrual"
	"github.com/vskurikhin/gophermart/internal/handlers"
	"github.com/vskurikhin/gophermart/internal/logger"
	"github.com/vskurikhin/gophermart/internal/model"
	"github.com/vskurikhin/gophermart/internal/storage"
	"github.com/vskurikhin/gophermart/internal/utils"
	"go.uber.org/zap"
	"net/http"
)

const numMsg = "number"

type number struct {
	log *zap.Logger
}

func newNumber() *number {
	return &number{log: logger.Get()}
}

// Handle number
//
//	@Summary		заказ
//	@Description	загрузка номера заказа
//	@Tags			Orders
//	@Accept			text/plain
//	@Produce		text/plain
//	@Param			request	body	string	true	"Формат запроса номер (body)"
//	@Security		BearerAuth
//	@Success		200				{string}	string			"номер заказа уже был загружен этим пользователем"
//	@Success		202				{string}	string			"новый номер заказа принят в обработку"
//	@Failure		400				{object}	model.JSONError	"неверный формат запроса"
//	@Failure		401				{object}	model.JSONError	"пользователь не аутентифицирован"
//	@Failure		409				{object}	model.JSONError	"номер заказа уже был загружен другим пользователем"
//	@Failure		422				{object}	model.JSONError	"неверный формат номера заказа"
//	@Failure		500				{string}	string			"Internal Server Error"
//	@Router			/user/orders 	[post]
//
//goland:noinspection GoUnhandledErrorResult
func (n *number) Handle(response http.ResponseWriter, request *http.Request) {

	ctx := request.Context()
	login, number, err := newRequestOrder(request).LoginNumber()

	if err != nil {

		n.log.Debug(numMsg, utils.LogCtxRecoverFields(ctx, err)...)
		render.Status(request, http.StatusBadRequest)
		render.Render(response, request, model.Error(handlers.ErrBadRequest))

		return
	}
	if login == nil {

		render.Status(request, http.StatusUnauthorized)
		render.Render(response, request, model.Error(handlers.ErrBadUserPassword))

		return
	}
	result := newService(ctx, storage.NewPgsStorage(), accrual.GetWorkers()).
		Number(*login, *number)
	render.Status(request, result.Status())

	switch value := result.(type) {
	case *handlers.ResultError:
		render.Render(response, request, model.Error(value.Error()))
	case *handlers.ResultString:
		render.Render(response, request, model.NewNumber(value.String()))
	default:
		n.log.Debug(numMsg, utils.InternalErrorZapField(ctx, request, result)...)
	}
}
