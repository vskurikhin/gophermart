/*
 * This file was last modified at 2024-04-15 10:39 by Victor N. Skurikhin.
 * number.go
 * $Id$
 */

package orders

import (
	"github.com/vskurikhin/gophermart/internal/logger"
	"go.uber.org/zap"
	"net/http"
)

type number struct {
	log *zap.Logger
}

func newNumber() *number {
	return &number{log: logger.Get()}
}

func (r *number) Handle(response http.ResponseWriter, request *http.Request) {

}
