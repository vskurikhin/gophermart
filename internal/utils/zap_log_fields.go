/*
 * This file was last modified at 2024-04-19 18:18 by Victor N. Skurikhin.
 * zap_log_fields.go
 * $Id$
 */

package utils

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/vskurikhin/gophermart/internal/handlers"
	"go.uber.org/zap"
	"net/http"
)

func LogCtxRecoverFields(ctx context.Context, r any) []zap.Field {
	return []zap.Field{
		zap.String("reqId", middleware.GetReqID(ctx)),
		zap.Reflect("error", r),
	}
}

func LogCtxReasonErrFields(ctx context.Context, reason string, err error) []zap.Field {
	return []zap.Field{
		zap.String("reqId", middleware.GetReqID(ctx)),
		zap.String("reason", reason),
		zap.Reflect("error", err.Error()),
	}
}

func InternalErrorZapField(ctx context.Context, request *http.Request, result handlers.Result) []zap.Field {
	err := fmt.Errorf("error result: %v", result)
	render.Status(request, http.StatusInternalServerError)
	return LogCtxReasonErrFields(ctx, err.Error(), handlers.ErrInternalError)
}
