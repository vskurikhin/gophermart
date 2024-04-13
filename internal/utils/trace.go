/*
 * This file was last modified at 2024-04-13 19:03 by Victor N. Skurikhin.
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

	l := logger.Get()

	if l.Level() != zap.DebugLevel {
		return func() {

		}
	}
	var uuid uuid4.UUID

	if value := ctx.Value("uuid"); value != nil {
		if uuidValue, ok := value.(uuid4.UUID); ok {
			uuid = uuidValue
		} else {
			uuid = uuid4.New()
		}
	} else {
		uuid = uuid4.New()
	}
	start := time.Now()

	f := fmt.Sprintf(" in[%s]: %s(", uuid, name) + format + ")"
	l.Debug(fmt.Sprintf(f, values...))

	return func() {
		l.Debug(fmt.Sprintf("out[%s]: %s [%s]", uuid, name, time.Since(start)))
	}
}
