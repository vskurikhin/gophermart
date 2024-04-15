/*
 * This file was last modified at 2024-04-15 10:39 by Victor N. Skurikhin.
 * withdrawals.go
 * $Id$
 */

package accounts

import (
	"github.com/vskurikhin/gophermart/internal/logger"
	"go.uber.org/zap"
	"net/http"
)

type withdrawals struct {
	log *zap.Logger
}

func newWithdraws() *withdrawals {
	return &withdrawals{log: logger.Get()}
}

func (r *withdrawals) Handle(response http.ResponseWriter, request *http.Request) {

}
