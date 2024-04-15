/*
 * This file was last modified at 2024-04-16 10:06 by Victor N. Skurikhin.
 * password.go
 * $Id$
 */

package utils

import "golang.org/x/crypto/bcrypt"

const (
	DefaultCost int = 9
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
