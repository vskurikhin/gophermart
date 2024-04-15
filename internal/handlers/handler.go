/*
 * This file was last modified at 2024-04-16 17:02 by Victor N. Skurikhin.
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
var ErrBadUserPassword = fmt.Errorf("неверная пара логин/пароль")
var ErrStatusConflict = fmt.Errorf("логин уже занят")
var ErrUserUnauthorized = fmt.Errorf("пользователь не аутентифицирован")
