/*
 * This file was last modified at 2024-04-15 10:39 by Victor N. Skurikhin.
 * handler.go
 * $Id$
 */

package handlers

import "net/http"

type Handler interface {
	Handle(response http.ResponseWriter, request *http.Request)
}
