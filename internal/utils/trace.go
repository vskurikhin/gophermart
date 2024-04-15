/*
 * This file was last modified at 2024-04-15 10:39 by Victor N. Skurikhin.
 * trace.go
 * $Id$
 */

package utils

import (
	"context"
	"fmt"
	uuid4 "github.com/google/uuid"
	"github.com/vskurikhin/gophermart/internal/logger"
	"go.uber.org/zap"
	"time"
)

func TraceInOut(ctx context.Context, name, format string, values ...any) func() {

	log := logger.Get()

	if log.Level() != zap.DebugLevel {
		return func() {}
	}
	uuid := GetUUID(ctx)
	start := time.Now()
	f := fmt.Sprintf(" in[%s]: %s(", uuid, name) + format + ")"
	log.Debug(fmt.Sprintf(f, values...))

	return func() {
		log.Debug(fmt.Sprintf("out[%s]: %s [%s]", uuid, name, time.Since(start)))
	}
}

func GetUUID(ctx context.Context) uuid4.UUID {

	if uuid, ok := ctx.Value("uuid").(uuid4.UUID); ok {
		return uuid
	}
	return uuid4.New()
}
