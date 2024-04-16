/*
 * This file was last modified at 2024-04-18 23:06 by Victor N. Skurikhin.
 * handle_result.go
 * $Id$
 */

package auth

import (
	"context"
	"fmt"
	"github.com/go-chi/render"
	"github.com/vskurikhin/gophermart/internal/handlers"
	"github.com/vskurikhin/gophermart/internal/logger"
	"github.com/vskurikhin/gophermart/internal/model"
	"github.com/vskurikhin/gophermart/internal/utils"
	"go.uber.org/zap"
	"net/http"
)

type handleResult struct {
	log      *zap.Logger
	response http.ResponseWriter
	request  *http.Request
}

func newHandleResult(response http.ResponseWriter, request *http.Request) *handleResult {
	return &handleResult{log: logger.Get(), response: response, request: request}
}

func (h *handleResult) handleUser(msg string, resultFunc func(ctx context.Context, user *model.User) handlers.Result) {

	ctx := h.request.Context()
	user, err := model.UnmarshalFromReader(h.request.Body)

	if err != nil || user == nil {

		h.log.Debug(msg, utils.LogCtxRecoverFields(ctx, err)...)
		render.Status(h.request, http.StatusBadRequest)
		//goland:noinspection GoUnhandledErrorResult
		render.Render(h.response, h.request, model.Error(handlers.ErrBadRequest))

		return
	}

	result := resultFunc(ctx, user)

	switch result.(type) {
	case *handlers.ResultError:

		render.Status(h.request, result.Status())
		e := result.(*handlers.ResultError)
		//goland:noinspection GoUnhandledErrorResult
		render.Render(h.response, h.request, model.Error(e.Error()))

	case *handlers.ResultString:

		t := result.(*handlers.ResultString)
		http.SetCookie(h.response, utils.NewCookie(t.String()))

		if err := user.MarshalToWriter(h.response); err != nil {
			panic(err)
		}
	default:

		err := fmt.Errorf("unknow result: %v", result)
		h.log.Debug(msg, utils.LogCtxRecoverFields(ctx, err)...)
		render.Status(h.request, http.StatusInternalServerError)
	}
}
