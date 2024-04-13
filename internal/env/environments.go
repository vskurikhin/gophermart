/*
 * This file was last modified at 2024-04-13 17:14 by Victor N. Skurikhin.
 * environments.go
 * $Id$
 */

package env

import (
	"fmt"
	"github.com/caarlos0/env"
	"strconv"
)

type environments struct {
	accrualSystemAddress []string `env:"ACCRUAL_SYSTEM_ADDRESS" envSeparator:":"`
	address              []string `env:"RUN_ADDRESS" envSeparator:":"`
	dataBaseDSN          string   `env:"DATABASE_DSN"`
	key                  string   `env:"KEY"`
}

func newEnvironments() *environments {

	e := new(environments)

	if err := env.Parse(e); err != nil {
		panic(err)
	}
	return e
}

func (e *environments) AccrualSystemAddress() string {
	return parseAddress(e.accrualSystemAddress)
}

func (e *environments) Address() string {
	return parseAddress(e.address)
}

func (e *environments) DataBaseDSN() string {
	return e.dataBaseDSN
}

func (e *environments) Key() string {
	return e.key
}

func parseAddress(address []string) string {

	if len(address) != 2 {
		return ""
	}
	port, err := strconv.Atoi(address[1])

	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s:%d", address[0], port)
}
