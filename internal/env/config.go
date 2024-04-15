/*
 * This file was last modified at 2024-04-16 09:51 by Victor N. Skurikhin.
 * config.go
 * $Id$
 */

package env

import (
	"sync"
)

type Config interface {
	AccrualSystemAddress() string
	Address() string
	DataBaseDSN() string
	DevelopmentLogger() bool
	Key() string
}

type config struct {
	accrualSystemAddress string
	address              string
	dataBaseDSN          string
	developmentLogger    bool
	key                  string
}

var once = new(sync.Once)
var cfg *config

func GetConfig() Config {

	once.Do(func() {
		cfg = new(config)
		e := newEnvironments()
		f := newFlags()

		if e.accrualSystemAddress() != "" {
			cfg.accrualSystemAddress = e.accrualSystemAddress()
		} else {
			cfg.accrualSystemAddress = *f.AccrualSystemAddress()
		}

		if e.address() != "" {
			cfg.address = e.address()
		} else {
			cfg.address = *f.Address()
		}

		if e.DataBaseDSN != "" {
			cfg.dataBaseDSN = e.DataBaseDSN
		} else {
			cfg.dataBaseDSN = *f.DataBaseDSN()
		}

		if e.Key != "" {
			cfg.key = e.Key
		} else {
			cfg.key = *f.Key()
		}
		cfg.developmentLogger = f.DevelopmentLogger()
	})

	return cfg
}

func (c *config) AccrualSystemAddress() string {
	return c.accrualSystemAddress
}

func (c *config) Address() string {
	return c.address
}

func (c *config) DataBaseDSN() string {
	return c.dataBaseDSN
}

func (c *config) DevelopmentLogger() bool {
	return c.developmentLogger
}

func (c *config) Key() string {
	return c.key
}
