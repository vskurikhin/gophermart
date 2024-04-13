/*
 * This file was last modified at 2024-04-13 17:14 by Victor N. Skurikhin.
 * statuses.go
 * $Id$
 */

package dao

import (
	"github.com/vskurikhin/gophermart/internal/domain/entity"
	"github.com/vskurikhin/gophermart/internal/storage"
)

type statuses struct {
	storage storage.Storage
}

func Statuses(storage storage.Storage) *statuses {
	return &statuses{storage: storage}
}

func (s *statuses) GetAllStatuses() ([]*entity.Status, error) {
	return entity.FuncGetAllStatuss()(s.storage)
}

func (s *statuses) GetStatus(id int) (*entity.Status, error) {
	return entity.FuncGetStatus()(s.storage, id)
}

func (s *statuses) Save(status *entity.Status) error {
	return status.Save(s.storage)
}
