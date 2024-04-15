/*
 * This file was last modified at 2024-04-15 10:39 by Victor N. Skurikhin.
 * balance.go
 * $Id$
 */

package accounts

import (
	"github.com/vskurikhin/gophermart/internal/logger"
	"go.uber.org/zap"
	"net/http"
)

type balance struct {
	log *zap.Logger
}

func newBalance() *balance {
	return &balance{log: logger.Get()}
}

func (r *balance) Handle(response http.ResponseWriter, request *http.Request) {

}
