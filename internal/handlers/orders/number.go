/*
 * This file was last modified at 2024-04-25 22:23 by Victor N. Skurikhin.
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
