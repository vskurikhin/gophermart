/*
 * This file was last modified at 2024-05-07 16:32 by Victor N. Skurikhin.
 * error.go
 * $Id$
 */

package model

import (
	"net/http"
)

type JSONError struct {
	Error string `json:"error"`
}

func Error(err error) *JSONError {
	return &JSONError{
		Error: err.Error(),
	}
}

func (e *JSONError) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
