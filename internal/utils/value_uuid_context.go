/*
 * This file was last modified at 2024-05-07 14:46 by Victor N. Skurikhin.
 * value_uuid_context.go
 * $Id$
 */

package utils

import (
	"context"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

func NewIDContext() context.Context {
	return context.WithValue(context.Background(), middleware.RequestIDKey, uuid.New().String())
}
