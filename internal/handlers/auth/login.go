/*
 * This file was last modified at 2024-04-15 10:39 by Victor N. Skurikhin.
 * login.go
 * $Id$
 */

package auth

import (
	"github.com/vskurikhin/gophermart/internal/logger"
	"go.uber.org/zap"
	"net/http"
)

type login struct {
	log *zap.Logger
}

func newLogin() *login {
	return &login{log: logger.Get()}
}

func (r *login) Handle(response http.ResponseWriter, request *http.Request) {

}
