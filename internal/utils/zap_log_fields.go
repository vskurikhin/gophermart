/*
 * This file was last modified at 2024-04-15 23:30 by Victor N. Skurikhin.
 * zap_log_fields.go
 * $Id$
 */

package utils

import (
	"context"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
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
		zap.Reflect("error", err),
	}
}
