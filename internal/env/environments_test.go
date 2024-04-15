/*
 * This file was last modified at 2024-04-15 22:15 by Victor N. Skurikhin.
 * environments_test.go
 * $Id$
 */

package env

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEnvironments(t *testing.T) {

	urlHost := "localhost:8080"

	var tests = []struct {
		name     string
		input    environments
		dontWant string
	}{
		{
			name: "Test environments positive #1",
			input: environments{
				AccrualSystemAddress: []string{},
				Address:              []string{"localhost", "8080"},
			},
			dontWant: "dontWant",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.input.accrualSystemAddress()
			assert.Equal(t, "", got)
			got = test.input.address()
			assert.Equal(t, urlHost, got)
			e := newEnvironments()
			got = e.DataBaseDSN
			assert.NotEqual(t, test.dontWant, got)
			got = e.Key
			assert.NotEqual(t, test.dontWant, got)
		})
	}
}
