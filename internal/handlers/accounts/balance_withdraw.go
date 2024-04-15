/*
 * This file was last modified at 2024-04-15 10:39 by Victor N. Skurikhin.
 * balance_withdraw.go
 * $Id$
 */

package accounts

import (
	"github.com/vskurikhin/gophermart/internal/logger"
	"go.uber.org/zap"
	"net/http"
)

type balanceWithdraw struct {
	log *zap.Logger
}

func newBalanceWithdraw() *balanceWithdraw {
	return &balanceWithdraw{log: logger.Get()}
}

func (r *balanceWithdraw) Handle(response http.ResponseWriter, request *http.Request) {

}
