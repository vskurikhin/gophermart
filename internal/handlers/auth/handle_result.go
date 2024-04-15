/*
 * This file was last modified at 2024-04-16 10:44 by Victor N. Skurikhin.
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

	if err, ok1 := result.(*handlers.ResultError); ok1 {

		render.Status(h.request, err.Status())
		//goland:noinspection GoUnhandledErrorResult
		render.Render(h.response, h.request, model.Error(err.Error()))

		return
	} else if token, ok2 := result.(*handlers.ResultString); ok2 {

		http.SetCookie(h.response, utils.NewCookie(token.String()))

		if err := user.MarshalToWriter(h.response); err != nil {
			panic(err)
		}
		return
	} else {
		err := fmt.Errorf("unknow result: %v, is result: %v, is error: %v", result, ok2, ok1)
		h.log.Debug(msg, utils.LogCtxRecoverFields(ctx, err)...)
	}
	render.Status(h.request, http.StatusInternalServerError)
}
