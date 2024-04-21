/*
 * This file was last modified at 2024-04-19 22:10 by Victor N. Skurikhin.
 * check_luhn_test.go
 * $Id$
 */

package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

//goland:noinspection GoImportUsedAsName,GoUnhandledErrorResult
func TestCheckLuhn(t *testing.T) {

	var tests = []struct {
		name   string
		number string
		result bool
	}{
		{"Test negative #1", "4561261212345464", false},
		{"Test positive #1", "4561261212345467", true},
		{"Test positive #2", "12345678903", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := CheckLuhn(test.number)
			assert.Equal(t, test.result, result)
		})
	}
}
