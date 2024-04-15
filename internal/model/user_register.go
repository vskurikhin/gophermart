/*
 * This file was last modified at 2024-04-15 14:11 by Victor N. Skurikhin.
 * user_register.go
 * $Id$
 */

package model

import (
	"github.com/mailru/easyjson"
	"io"
)

type UserRegister struct {
	Login    string `json:"login"`    // логин
	Password string `json:"password"` // пароль
}

func UnmarshalFromReader(reader io.Reader) (*UserRegister, error) {

	userRegister := new(UserRegister)

	if err := easyjson.UnmarshalFromReader(reader, userRegister); err != nil {
		return nil, err
	}
	return userRegister, nil
}

func (u *UserRegister) MarshalToWriter(writer io.Writer) error {

	if _, err := easyjson.MarshalToWriter(u, writer); err != nil {
		return err
	}
	return nil
}
