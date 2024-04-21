/*
 * This file was last modified at 2024-04-20 18:04 by Victor N. Skurikhin.
 * handle_result.go
 * $Id$
 */

package auth

import (
	"context"
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

//goland:noinspection GoUnhandledErrorResult
func (h *handleResult) handleUser(msg string, resultFunc func(ctx context.Context, user *model.User) handlers.Result) {

	ctx := h.request.Context()
	user, err := model.UnmarshalUserFromReader(h.request.Body)

	if err != nil || user == nil {

		h.log.Debug(msg, utils.LogCtxRecoverFields(ctx, err)...)
		render.Status(h.request, http.StatusBadRequest)
		render.Render(h.response, h.request, model.Error(handlers.ErrBadRequest))

		return
	}

	result := resultFunc(ctx, user)
	render.Status(h.request, result.Status())

	switch value := result.(type) {
	case *handlers.ResultError:
		render.Render(h.response, h.request, model.Error(value.Error()))
	case *handlers.ResultString:

		http.SetCookie(h.response, utils.NewCookie(value.String()))

		h.response.Header().Set("Content-Type", "application/json")
		if err := user.MarshalToWriter(h.response); err == nil {
			return
		}
	}
	h.log.Debug(msg, utils.InternalErrorZapField(ctx, h.request, result)...)
}
