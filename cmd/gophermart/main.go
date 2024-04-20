/*
 * This file was last modified at 2024-04-20 00:07 by Victor N. Skurikhin.
 * main.go
 * $Id$
 */

package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/vskurikhin/gophermart/internal/env"
	"github.com/vskurikhin/gophermart/internal/handlers/accounts"
	"github.com/vskurikhin/gophermart/internal/handlers/auth"
	"github.com/vskurikhin/gophermart/internal/handlers/orders"
	"github.com/vskurikhin/gophermart/internal/logger"
	"github.com/vskurikhin/gophermart/internal/utils"
	"github.com/vskurikhin/gophermart/internal/zapchi"
	"net/http"
)

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(zapchi.Logger(logger.Get(), "router"))
	router.Use(render.SetContentType(render.ContentTypeJSON))

	router.Group(func(r chi.Router) {
		r.Post("/api/user/login", auth.UserLoginHandlerFunc())
		r.Post("/api/user/register", auth.UserRegisterHandlerFunc())
	})

	router.Group(func(r chi.Router) {
		r.Use(utils.Verifier())
		r.Use(utils.UnauthenticatedError)
		r.Post("/api/user/orders", orders.UserNumberHandlerFunc())
	})

	router.Group(func(r chi.Router) {
		r.Use(utils.Verifier())
		r.Use(utils.UnauthorizedError)
		r.Get("/api/user/orders", orders.UserOrdersHandlerFunc())

		r.Get("/api/user/balance", accounts.BalanceHandlerFunc())
		r.Post("/api/user/balance/withdraw", accounts.BalanceWithdrawHandlerFunc())
		r.Get("/api/user/withdrawals", accounts.WithdrawalsHandlerFunc())
	})

	err := http.ListenAndServe(env.GetConfig().Address(), router)
	if err != nil {
		panic(err)
	}
}
