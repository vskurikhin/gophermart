package dao

import (
	"github.com/vskurikhin/gophermart/internal/domain/entity"
	"github.com/vskurikhin/gophermart/internal/storage"
)

type status struct {
	storage storage.Storage
}

func Statuses(storage storage.Storage) *status {
	return &status{storage: storage}
}

func (w *status) GetAllStatuses() ([]*entity.Status, error) {
	return entity.FuncGetAllStatuss()(w.storage)
}

func (w *status) GetStatus(id int) (*entity.Status, error) {
	return entity.FuncGetStatus()(w.storage, id)
}

func (w *status) Save(status *entity.Status) error {
	return status.Save(w.storage)
}
