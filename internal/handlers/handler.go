/*
 * This file was last modified at 2024-04-15 16:54 by Victor N. Skurikhin.
 * handler.go
 * $Id$
 */

package handlers

import (
	"fmt"
	"net/http"
)

type Handler interface {
	Handle(response http.ResponseWriter, request *http.Request)
}

var ErrBadRequest = fmt.Errorf("неверный формат запроса")
var ErrStatusConflict = fmt.Errorf("логин уже занят")
