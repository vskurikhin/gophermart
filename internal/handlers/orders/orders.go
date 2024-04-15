/*
 * This file was last modified at 2024-04-15 10:39 by Victor N. Skurikhin.
 * orders.go
 * $Id$
 */

package orders

import (
	"github.com/vskurikhin/gophermart/internal/logger"
	"go.uber.org/zap"
	"net/http"
)

type orders struct {
	log *zap.Logger
}

func newOrders() *orders {
	return &orders{log: logger.Get()}
}

func (r *orders) Handle(response http.ResponseWriter, request *http.Request) {

}
