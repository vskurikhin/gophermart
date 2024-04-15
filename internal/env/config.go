/*
 * This file was last modified at 2024-04-15 11:23 by Victor N. Skurikhin.
 * config.go
 * $Id$
 */

package env

import (
	"github.com/vskurikhin/gophermart/internal/logger"
	"sync"
)

type Config interface {
	AccrualSystemAddress() string
	Address() string
	DataBaseDSN() string
	Key() string
}

type config struct {
	accrualSystemAddress string
	address              string
	dataBaseDSN          string
	key                  string
}

var once = new(sync.Once)
var cfg *config

func GetConfig() Config {

	once.Do(func() {
		cfg = new(config)
		e := newEnvironments()
		f := newFlags()

		if e.AccrualSystemAddress() != "" {
			cfg.accrualSystemAddress = e.AccrualSystemAddress()
		} else {
			cfg.accrualSystemAddress = *f.AccrualSystemAddress()
		}

		if e.Address() != "" {
			cfg.address = e.Address()
		} else {
			cfg.address = *f.Address()
		}

		if e.DataBaseDSN() != "" {
			cfg.dataBaseDSN = e.DataBaseDSN()
		} else {
			cfg.dataBaseDSN = *f.DataBaseDSN()
		}

		if e.Key() != "" {
			cfg.key = e.Key()
		} else {
			cfg.key = *f.Key()
		}
		logger.Development = *f.DevelopmentLogger()
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

func (c *config) Key() string {
	return c.key
}
