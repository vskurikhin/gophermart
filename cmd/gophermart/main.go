/*
 * This file was last modified at 2024-04-15 14:48 by Victor N. Skurikhin.
 * main.go
 * $Id$
 */

package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/vskurikhin/gophermart/internal/env"
	"github.com/vskurikhin/gophermart/internal/handlers/auth"
	"github.com/vskurikhin/gophermart/internal/storage"
	"net/http"
)

func main() {
	_ = storage.GetDB()
	router := chi.NewRouter()
	router.Use(render.SetContentType(render.ContentTypeJSON))
	router.Post("/api/user/login", auth.UserLoginHandlerFunc())
	router.Post("/api/user/register", auth.UserRegisterHandlerFunc())

	err := http.ListenAndServe(env.GetConfig().Address(), router)
	if err != nil {
		panic(err)
	}
}
