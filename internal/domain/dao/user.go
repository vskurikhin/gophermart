package dao

import (
	"github.com/vskurikhin/gophermart/internal/domain/entity"
	"github.com/vskurikhin/gophermart/internal/storage"
)

type user struct {
	storage storage.Storage
}

func Users(storage storage.Storage) *user {
	return &user{storage: storage}
}

func (w *user) GetAllUsers() ([]*entity.User, error) {
	return entity.FuncGetAllUsers()(w.storage)
}

func (w *user) GetUser(login string) (*entity.User, error) {
	return entity.FuncGetUser()(w.storage, login)
}

func (w *user) Save(user *entity.User) (*entity.User, error) {
	return user.Save(w.storage)
}
