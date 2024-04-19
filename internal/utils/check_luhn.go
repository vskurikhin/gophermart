/*
 * This file was last modified at 2024-04-19 23:24 by Victor N. Skurikhin.
 * check_luhn.go
 * $Id$
 */

package utils

func CheckLuhn(number string) bool {

	if len(number) < 1 {
		return false
	}

	var sum int
	parity := len(number) % 2

	for i := 0; i < len(number); i++ {
		digit := int(number[i] - '0')

		if i%2 == parity {
			digit *= 2

			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}
	return sum%10 == 0
}
