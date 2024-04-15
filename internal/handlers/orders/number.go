/*
 * This file was last modified at 2024-04-16 12:41 by Victor N. Skurikhin.
 * number.go
 * $Id$
 */

package orders

import (
	"fmt"
	"github.com/go-chi/jwtauth"
	"github.com/vskurikhin/gophermart/internal/logger"
	"go.uber.org/zap"
	"net/http"
	"os"
)

type number struct {
	log *zap.Logger
}

func newNumber() *number {
	return &number{log: logger.Get()}
}

func (r *number) Handle(response http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	_, m, _ := jwtauth.FromContext(ctx)
	fmt.Fprintf(os.Stderr, "Ok number. username: %+v\n", m["username"])
}
