/*
 * This file was last modified at 2024-04-15 14:12 by Victor N. Skurikhin.
 * error.go
 * $Id$
 */

package model

import (
	"net/http"
)

type rError struct {
	Error string `json:"error"`
}

func Error(err error) *rError {
	return &rError{
		Error: err.Error(),
	}
}

func (e *rError) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
