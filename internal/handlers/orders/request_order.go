/*
 * This file was last modified at 2024-04-19 14:39 by Victor N. Skurikhin.
 * request_order.go
 * $Id$
 */

package orders

import (
	"context"
	"github.com/vskurikhin/gophermart/internal/handlers"
	"github.com/vskurikhin/gophermart/internal/handlers/auth"
	"io"
	"log"
	"net/http"
	"regexp"
)

type RequestOrder interface {
	Context() context.Context
}

type requestOrder struct {
	request RequestOrder
}

func newRequestOrder(request RequestOrder) *requestOrder {
	return &requestOrder{request: request}
}

func (o requestOrder) LoginNumber() (*string, *string, error) {

	login, err := o.Login()

	if err != nil {
		return nil, nil, err
	}

	if request, ok := o.request.(*http.Request); ok {

		number, err := Number(request)

		if err != nil {
			return nil, nil, err
		}
		return login, number, err
	}

	return nil, nil, handlers.ErrBadRequest
}

func (o *requestOrder) Login() (*string, error) {
	ctx := o.request.Context()
	return auth.Login(ctx)
}

var reTextPlain, _ = regexp.Compile(`^text/plain.*`)

func Number(request *http.Request) (*string, error) {

	contentType := request.Header.Get("Content-Type")

	if reTextPlain.MatchString(contentType) {

		b, err := io.ReadAll(request.Body)
		if err != nil {
			log.Fatalln(err)
		}
		number := string(b)

		return &number, nil
	}
	return nil, handlers.ErrBadRequest
}
