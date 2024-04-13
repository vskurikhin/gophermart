/*
 * This file was last modified at 2024-04-12 14:53 by Victor N. Skurikhin.
 * trace.go
 * $Id$
 */

package util

import (
	"context"
	"fmt"
	uuid4 "github.com/google/uuid"
	"log"
	"time"
)

func Trace(ctx context.Context, name, format string, values ...any) func() {

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
	log.Printf(f, values...)

	return func() {
		log.Printf("out[%s]: %s [%s]", uuid, name, time.Since(start))
	}
}
