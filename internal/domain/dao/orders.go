/*
 * This file was last modified at 2024-04-16 18:36 by Victor N. Skurikhin.
 * orders.go
 * $Id$
 */

package dao

import (
	"github.com/vskurikhin/gophermart/internal/domain/entity"
	"github.com/vskurikhin/gophermart/internal/storage"
)

type orders struct {
	storage storage.Storage
}

func Orders(storage storage.Storage) *orders {
	return &orders{storage: storage}
}

func (o *orders) GetAllOrders() ([]*entity.Order, error) {
	return entity.FuncGetAllOrders()(o.storage)
}

func (o *orders) GetOrder(login, number string) (*entity.Order, error) {
	return entity.FuncGetOrder()(o.storage, login, number)
}

func (o *orders) Insert(order *entity.Order) (*entity.Order, error) {
	return order.Insert(o.storage)
}

func (o *orders) Save(order *entity.Order) (*entity.Order, error) {
	return order.Save(o.storage)
}
