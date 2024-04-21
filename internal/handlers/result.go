/*
 * This file was last modified at 2024-04-16 09:38 by Victor N. Skurikhin.
 * result.go
 * $Id$
 */

package handlers

type Result interface {
	Result() any
	Status() int
}

type ResultAny struct {
	result any
	status int
}

func NewResultAny(result any, status int) *ResultAny {
	return &ResultAny{result: result, status: status}
}

func (r *ResultAny) Result() any {
	return r.result
}

func (r *ResultAny) Status() int {
	return r.status
}

type ResultError struct {
	err    error
	status int
}

func NewResultError(error error, status int) *ResultError {
	return &ResultError{err: error, status: status}
}

func (r *ResultError) Error() error {
	return r.err
}

func (r *ResultError) Result() any {
	return r.err
}

func (r *ResultError) Status() int {
	return r.status
}

type ResultString struct {
	result string
	status int
}

func NewResultString(result string, status int) *ResultString {
	return &ResultString{result: result, status: status}
}

func (r *ResultString) Result() any {
	return r.result
}

func (r *ResultString) Status() int {
	return r.status
}

func (r *ResultString) String() string {
	return r.result
}
