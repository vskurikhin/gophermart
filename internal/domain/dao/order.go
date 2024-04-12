/*
 * This file was last modified at 2024-04-12 16:30 by Victor N. Skurikhin.
 * order.go
 * $Id$
 */

package dao

import (
	"github.com/vskurikhin/gophermart/internal/domain/entity"
	"github.com/vskurikhin/gophermart/internal/storage"
)

type order struct {
	storage storage.Storage
}

func Orders(storage storage.Storage) *order {
	return &order{storage: storage}
}

func (w *order) GetAllOrders() ([]*entity.Order, error) {
	return entity.FuncGetAllOrders()(w.storage)
}

func (w *order) GetOrder(login, number string) (*entity.Order, error) {
	return entity.FuncGetOrder()(w.storage, login, number)
}

func (w *order) Save(order *entity.Order) (*entity.Order, error) {
	return order.Save(w.storage)
}
