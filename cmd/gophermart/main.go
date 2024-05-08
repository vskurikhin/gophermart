/*
 * This file was last modified at 2024-05-07 18:50 by Victor N. Skurikhin.
 * main.go
 * $Id$
 */

package main

import (
	"embed"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/vskurikhin/gophermart/docs"
	"github.com/vskurikhin/gophermart/internal/env"
	"github.com/vskurikhin/gophermart/internal/handlers/accounts"
	"github.com/vskurikhin/gophermart/internal/handlers/auth"
	"github.com/vskurikhin/gophermart/internal/handlers/orders"
	"github.com/vskurikhin/gophermart/internal/logger"
	"github.com/vskurikhin/gophermart/internal/storage"
	"github.com/vskurikhin/gophermart/internal/utils"
	"github.com/vskurikhin/gophermart/internal/zapchi"
	"log"
	"net/http"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS //

//	@title			GopherMart API
//	@version		1.0
//	@description	This is a sample server GopherMart server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/api

//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization

func main() {

	dbMigrations()
	router := routerSetup()
	err := http.ListenAndServe(env.GetConfig().Address(), router)

	if err != nil {
		panic(err)
	}
}

func dbMigrations() {
	pool, ok := storage.GetDB().DBPool()
	if !ok {
		panic(fmt.Errorf("pool: %v not ok", pool))
	}

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	db := stdlib.OpenDBFromPool(pool)
	if err := goose.Up(db, "migrations"); err != nil {
		panic(err)
	}
	if err := goose.Version(db, "migrations"); err != nil {
		log.Fatal(err)
	}
	if err := db.Close(); err != nil {
		panic(err)
	}
}

func routerSetup() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(zapchi.Logger(logger.Get(), "router"))
	router.Use(render.SetContentType(render.ContentTypeJSON))

	router.Group(func(r chi.Router) {
		r.Use(middleware.Compress(9))
		r.Post("/api/user/login", auth.UserLoginHandlerFunc())
		r.Post("/api/user/register", auth.UserRegisterHandlerFunc())
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"), //The url pointing to API definition
		))
	})

	//	@Security	BearerAuth
	router.Group(func(r chi.Router) {
		r.Use(utils.Verifier())
		r.Use(utils.UnauthenticatedError)
		r.Post("/api/user/orders", orders.UserNumberHandlerFunc())
	})

	router.Group(func(r chi.Router) {
		r.Use(utils.Verifier())
		r.Use(utils.UnauthorizedError)
		r.Use(middleware.Compress(6))
		r.Get("/api/user/orders", orders.UserOrdersHandlerFunc())

		r.Get("/api/user/balance", accounts.BalanceHandlerFunc())
		r.Post("/api/user/balance/withdraw", accounts.BalanceWithdrawHandlerFunc())
		r.Get("/api/user/withdrawals", accounts.WithdrawalsHandlerFunc())
	})
	return router
}
