/*
 * This file was last modified at 2024-04-15 10:39 by Victor N. Skurikhin.
 * register.go
 * $Id$
 */

package auth

import (
	"github.com/vskurikhin/gophermart/internal/logger"
	"go.uber.org/zap"
	"net/http"
)

type register struct {
	log *zap.Logger
}

func newRegister() *register {
	return &register{log: logger.Get()}
}

func (r *register) Handle(response http.ResponseWriter, request *http.Request) {

}
