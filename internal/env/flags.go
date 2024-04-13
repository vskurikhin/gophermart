/*
 * This file was last modified at 2024-04-13 17:14 by Victor N. Skurikhin.
 * flags.go
 * $Id$
 */

package env

import "github.com/spf13/pflag"

type flags struct {
	accrualSystemAddress *string
	address              *string
	dataBaseDSN          *string
	key                  *string
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
		"postgresql://postgres:postgres@localhost/praktikum?sslmode=disable",
		"help message for file database DSN",
	)
	f.key = pflag.StringP(
		"key",
		"k",
		"",
		"help message for key",
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
