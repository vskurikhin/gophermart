/*
 * This file was last modified at 2024-04-13 19:08 by Victor N. Skurikhin.
 * config.go
 * $Id$
 */

package env

import (
	"github.com/vskurikhin/gophermart/internal/logger"
	"sync"
)

type Config struct {
	accrualSystemAddress string
	address              string
	dataBaseDSN          string
	key                  string
}

var once = new(sync.Once)
var config *Config

func GetConfig() *Config {

	once.Do(func() {
		config = new(Config)
		e := newEnvironments()
		f := newFlags()

		if e.AccrualSystemAddress() != "" {
			config.accrualSystemAddress = e.AccrualSystemAddress()
		} else {
			config.accrualSystemAddress = *f.AccrualSystemAddress()
		}

		if e.Address() != "" {
			config.address = e.Address()
		} else {
			config.address = *f.Address()
		}

		if e.DataBaseDSN() != "" {
			config.dataBaseDSN = e.DataBaseDSN()
		} else {
			config.dataBaseDSN = *f.DataBaseDSN()
		}

		if e.Key() != "" {
			config.key = e.Key()
		} else {
			config.key = *f.Key()
		}
		logger.Development = *f.DevelopmentLogger()
	})

	return config
}

func (c *Config) AccrualSystemAddress() string {
	return c.accrualSystemAddress
}

func (c *Config) Address() string {
	return c.address
}

func (c *Config) DataBaseDSN() string {
	return c.dataBaseDSN
}

func (c *Config) Key() string {
	return c.key
}
