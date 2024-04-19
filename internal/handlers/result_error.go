/*
 * This file was last modified at 2024-04-19 17:54 by Victor N. Skurikhin.
 * result_error.go
 * $Id$
 */

package handlers

import (
	"fmt"
	"net/http"
)

var ErrBadFormatNumber = fmt.Errorf("неверный формат номера заказа")
var ErrBadRequest = fmt.Errorf("неверный формат запроса")
var ErrBadUserPassword = fmt.Errorf("неверная пара логин/пароль")
var ErrBalanceNotSet = fmt.Errorf("баланс не установлен")
var ErrInternalError = fmt.Errorf("внутренняя ошибка сервера")
var ErrOrderByUserAlreadyLoaded = fmt.Errorf("номер заказа уже был загружен этим пользователем")
var ErrOrderOtherAlreadyLoaded = fmt.Errorf("номер заказа уже был загружен другим пользователем")
var ErrStatusConflict = fmt.Errorf("логин уже занят")
var ErrUserUnauthorized = fmt.Errorf("пользователь не аутентифицирован")

func ResultErrorBadFormatNumber() *ResultError {
	return NewResultError(ErrBadFormatNumber, http.StatusUnprocessableEntity)
}

func ResultErrorBadRequest() *ResultError {
	return NewResultError(ErrBadRequest, http.StatusBadRequest)
}

func ResultInternalError() *ResultError {
	return NewResultError(ErrInternalError, http.StatusInternalServerError)
}

func ResultErrorOrderByUserAlreadyLoaded() *ResultError {
	return NewResultError(ErrOrderByUserAlreadyLoaded, http.StatusOK)
}

func ResultErrorOrderOtherAlreadyLoaded() *ResultError {
	return NewResultError(ErrOrderOtherAlreadyLoaded, http.StatusConflict)
}
func ResultErrorStatusConflict() *ResultError {
	return NewResultError(ErrStatusConflict, http.StatusConflict)
}
