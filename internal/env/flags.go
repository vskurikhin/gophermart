/*
 * This file was last modified at 2024-04-20 17:09 by Victor N. Skurikhin.
 * flags.go
 * $Id$
 */

package env

import "github.com/spf13/pflag"

const DefaultConnDatabaseDSN = "postgresql://postgres:postgres@localhost/praktikum?sslmode=disable"

type flags struct {
	accrualSystemAddress *string
	address              *string
	dataBaseDSN          *string
	developmentLogger    *bool
	key                  *string
	rateLimit            *int
}

func (f *flags) RateLimit() int {
	if f.rateLimit == nil {
		return 60
	} else {
		return *f.rateLimit
	}
}

func newFlags() *flags {

	f := new(flags)
	f.accrualSystemAddress = pflag.StringP(
		"accrual-system-address",
		"r", "localhost:8079",
		"help message for accrual system address",
	)
	f.address = pflag.StringP(
		"run-address",
		"a",
		"localhost:8080",
		"help message for accrual run address",
	)
	f.dataBaseDSN = pflag.StringP(
		"database-dsn",
		"d",
		DefaultConnDatabaseDSN,
		"help message for file database DSN",
	)
	f.key = pflag.StringP(
		"key",
		"k",
		"",
		"help message for key",
	)
	f.developmentLogger = pflag.BoolP(
		"development-logger",
		"z",
		true,
		"help message for development logger",
	)
	f.rateLimit = pflag.IntP(
		"rate-limit",
		"l",
		60,
		"help message for rate limit",
	)
	pflag.Parse()

	return f
}

func (f *flags) AccrualSystemAddress() *string {
	return f.accrualSystemAddress
}

func (f *flags) Address() *string {
	return f.address
}

func (f *flags) DataBaseDSN() *string {
	return f.dataBaseDSN
}

func (f *flags) Key() *string {
	return f.key
}

func (f *flags) DevelopmentLogger() bool {
	return f.developmentLogger != nil && *f.developmentLogger
}
