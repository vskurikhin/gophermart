/*
 * This file was last modified at 2024-04-18 23:06 by Victor N. Skurikhin.
 * number.go
 * $Id$
 */

package orders

import (
	"fmt"
	"github.com/go-chi/render"
	"github.com/vskurikhin/gophermart/internal/handlers"
	"github.com/vskurikhin/gophermart/internal/logger"
	"github.com/vskurikhin/gophermart/internal/model"
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

func (r *number) Handle(response http.ResponseWriter, request *http.Request) {

	ctx := request.Context()
	login, number, err := newRequestOrder(request).LoginNumber()

	if err != nil || login == nil {

		r.log.Debug(numMsg, utils.LogCtxRecoverFields(ctx, err)...)
		render.Status(request, http.StatusBadRequest)
		//goland:noinspection GoUnhandledErrorResult
		render.Render(response, request, model.Error(handlers.ErrBadRequest))

		return
	}
	result := newService(ctx).Number(*login, *number)

	switch result.(type) {
	case *handlers.ResultError:

		render.Status(request, result.Status())
		e := result.(*handlers.ResultError)
		//goland:noinspection GoUnhandledErrorResult
		render.Render(response, request, model.Error(e.Error()))

	case *handlers.ResultString:

		render.Status(request, result.Status())
		t := result.(*handlers.ResultString)
		//goland:noinspection GoUnhandledErrorResult
		render.Render(response, request, model.NewNumber(t.String()))

	default:

		err := fmt.Errorf("unknow result: %v", result)
		r.log.Debug(numMsg, utils.LogCtxRecoverFields(ctx, err)...)
		render.Status(request, http.StatusInternalServerError)
	}
}
