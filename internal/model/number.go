/*
 * This file was last modified at 2024-04-16 18:42 by Victor N. Skurikhin.
 * number.go
 * $Id$
 */

package model

import "net/http"

type Number struct {
	Number string `json:"number"`
}

func NewNumber(number string) *Number {
	return &Number{
		Number: number,
	}
}

func (e *Number) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
