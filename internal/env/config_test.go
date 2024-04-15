/*
 * This file was last modified at 2024-04-15 11:32 by Victor N. Skurikhin.
 * config_test.go
 * $Id$
 */

package env

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig(t *testing.T) {

	var tests = []struct {
		name     string
		dontWant string
	}{
		{
			name:     "Test config positive #1",
			dontWant: " ",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := GetConfig()
			got := c.Address()
			assert.NotEqual(t, test.dontWant, got)
			got = c.AccrualSystemAddress()
			assert.NotEqual(t, test.dontWant, got)
			got = c.DataBaseDSN()
			assert.NotEqual(t, test.dontWant, got)
			got = c.Key()
			assert.NotEqual(t, test.dontWant, got)
		})
	}
}
