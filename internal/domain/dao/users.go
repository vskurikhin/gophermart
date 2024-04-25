/*
 * This file was last modified at 2024-04-25 21:54 by Victor N. Skurikhin.
 * users.go
 * $Id$
 */

package dao

import (
	"github.com/vskurikhin/gophermart/internal/domain/entity"
	"github.com/vskurikhin/gophermart/internal/storage"
)

type users struct {
	storage storage.Storage
}

func Users(storage storage.Storage) *users {
	return &users{storage: storage}
}

func (u *users) GetUser(login string) (*entity.User, error) {
	return entity.GetUser(u.storage, login)
}
