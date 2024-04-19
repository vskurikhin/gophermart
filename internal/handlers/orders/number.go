/*
 * This file was last modified at 2024-04-19 20:15 by Victor N. Skurikhin.
 * number.go
 * $Id$
 */

package orders

import (
	"github.com/go-chi/render"
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

func (n *number) Handle(response http.ResponseWriter, request *http.Request) {

	ctx := request.Context()
	login, number, err := newRequestOrder(request).LoginNumber()

	if err != nil || login == nil {

		n.log.Debug(numMsg, utils.LogCtxRecoverFields(ctx, err)...)
		render.Status(request, http.StatusBadRequest)
		//goland:noinspection GoUnhandledErrorResult
		render.Render(response, request, model.Error(handlers.ErrBadRequest))

		return
	}
	result := newService(ctx, storage.NewPgsStorage()).Number(*login, *number)

	switch value := result.(type) {
	case *handlers.ResultError:

		render.Status(request, result.Status())
		//goland:noinspection GoUnhandledErrorResult
		render.Render(response, request, model.Error(value.Error()))

	case *handlers.ResultString:

		render.Status(request, result.Status())
		//goland:noinspection GoUnhandledErrorResult
		render.Render(response, request, model.NewNumber(value.String()))

	default:
		n.log.Debug(numMsg, utils.InternalErrorZapField(ctx, request, result)...)
	}
}
