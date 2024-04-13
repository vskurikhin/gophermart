/*
 * This file was last modified at 2024-04-13 17:14 by Victor N. Skurikhin.
 * flags_test.go
 * $Id$
 */

package env

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFlags(t *testing.T) {

	urlHost := "localhost:8080"

	var tests = []struct {
		name            string
		input           flags
		want            string
		wantDataBaseDSN string
	}{
		{
			name: "Test flags positive #1",
			input: flags{
				accrualSystemAddress: &urlHost,
				address:              &urlHost,
			},
			want:            "localhost:8080",
			wantDataBaseDSN: "postgresql://postgres:postgres@localhost/praktikum?sslmode=disable",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.input.Address()
			assert.Equal(t, test.want, *got)
			got = test.input.AccrualSystemAddress()
			assert.Equal(t, test.want, *got)
		})
	}
}
