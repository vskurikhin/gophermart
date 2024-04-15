/*
 * This file was last modified at 2024-04-15 16:56 by Victor N. Skurikhin.
 * value_uuid_context.go
 * $Id$
 */

package utils

import (
	"context"
	"github.com/google/uuid"
)

type tUUIDKey string

const UUIDKey = tUUIDKey("uuid")

func NewUUIDContext() context.Context {
	return context.WithValue(context.Background(), UUIDKey, uuid.New())
}
