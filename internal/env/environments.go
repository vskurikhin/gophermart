/*
 * This file was last modified at 2024-04-20 17:09 by Victor N. Skurikhin.
 * environments.go
 * $Id$
 */

package env

import (
	"fmt"
	"github.com/caarlos0/env"
	"regexp"
	"strconv"
)

type environments struct {
	AccrualSystemAddress []string `env:"ACCRUAL_SYSTEM_ADDRESS" envSeparator:":"`
	Address              []string `env:"RUN_ADDRESS" envSeparator:":"`
	DataBaseDSN          string   `env:"DATABASE_URI"`
	Key                  string   `env:"KEY"`
}

func newEnvironments() *environments {

	e := new(environments)

	if err := env.Parse(e); err != nil {
		panic(err)
	}
	return e
}

func (e *environments) accrualSystemAddress() string {
	return parseAddress(e.AccrualSystemAddress)
}

func (e *environments) address() string {
	return parseAddress(e.Address)
}

var reHTTP, _ = regexp.Compile(`^http.*`)

func parseAddress(address []string) string {

	if len(address) < 2 {
		return ""
	}
	var host string

	if reHTTP.MatchString(address[0]) {
		address = address[1:]
		host = address[0][2:]
	} else if !reHTTP.MatchString(address[0]) {
		host = address[0]
	} else {
		return ""
	}
	port, err := strconv.Atoi(address[1])

	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s:%d", host, port)
}
