/*
 * This file was last modified at 2024-04-19 11:27 by Victor N. Skurikhin.
 * trace.go
 * $Id$
 */

package utils

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vskurikhin/gophermart/internal/logger"
	"go.uber.org/zap"
	"time"
)

func TraceInOut(ctx context.Context, name, format string, values ...any) func() {

	log := logger.Get()

	if log.Level() > zap.DebugLevel {
		return func() {}
	}
	reqID := middleware.GetReqID(ctx)
	start := time.Now()
	f := fmt.Sprintf(" in[%s]: %s(", reqID, name) + format + ")"
	log.Debug(fmt.Sprintf(f, values...))

	return func() {
		log.Debug(fmt.Sprintf("out[%s]: %s [%s]", reqID, name, time.Since(start)))
	}
}
