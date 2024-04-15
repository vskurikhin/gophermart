/*
 * This file was last modified at 2024-04-16 10:37 by Victor N. Skurikhin.
 * user.go
 * $Id$
 */

package model

import (
	"github.com/mailru/easyjson"
	"io"
)

type User struct {
	Login    string `json:"login"`    // логин
	Password string `json:"password"` // пароль
}

func UnmarshalFromReader(reader io.Reader) (*User, error) {

	userRegister := new(User)

	if err := easyjson.UnmarshalFromReader(reader, userRegister); err != nil {
		return nil, err
	}
	return userRegister, nil
}

func (u *User) MarshalToWriter(writer io.Writer) error {

	if _, err := easyjson.MarshalToWriter(u, writer); err != nil {
		return err
	}
	return nil
}
