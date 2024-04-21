/*
 * This file was last modified at 2024-04-15 23:01 by Victor N. Skurikhin.
 * cookie.go
 * $Id$
 */

package utils

import (
	"net/http"
	"time"
)

func NewCookie(token string) *http.Cookie {
	return &http.Cookie{
		HttpOnly: true,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		SameSite: http.SameSiteLaxMode,
		Name:     "jwt",
		Value:    token,
	}
}
