/*
 * This file was last modified at 2024-05-07 12:48 by Victor N. Skurikhin.
 * order.go
 * $Id$
 */

package view

import (
	"database/sql"
	"github.com/jackc/pgx/v5"
	"github.com/vskurikhin/gophermart/internal/storage"
	"math/big"
	"time"
)

type Order struct {
	login      string
	number     string
	status     string
	accrual    *big.Float
	uploadedAt *time.Time
	createdAt  time.Time
	updateAt   *time.Time
}

func GetAllOrdersForLogin(s storage.Storage, login string) ([]*Order, error) {

	rows, err := s.GetAllForString(
		`SELECT o.login, o.number, o.accrual, o.uploaded_at, o.created_at, o.update_at, s.status
			FROM "orders" o JOIN status s ON o.status_id = s.id  WHERE login = $1`, login,
	)
	if err != nil {
		return nil, err
	}
	result := make([]*Order, 0)
	for rows.Next() {
		order, err := extractOrder(rows)
		if err != nil {
			return result, err
		}
		result = append(result, order)
	}
	return result, nil
}

func (o *Order) Number() string {
	return o.number
}

func (o *Order) Status() string {
	return o.status
}

func (o *Order) Accrual() *big.Float {
	return o.accrual
}

func (o *Order) UploadedAt() *time.Time {
	return o.uploadedAt
}

func extractOrder(row pgx.Row) (*Order, error) {

	var login, number, status string
	var createdAt time.Time
	var accrualNullString sql.NullString
	var uploadedAtNullTime, updateAtNullTime sql.NullTime

	err := row.Scan(
		&login, &number, &accrualNullString, &uploadedAtNullTime, &createdAt, &updateAtNullTime, &status,
	)

	if err != nil {
		return nil, err
	}
	var uploadedAt *time.Time

	if uploadedAtNullTime.Valid {
		uploadedAt = &uploadedAtNullTime.Time
	}
	var updateAt *time.Time

	if updateAtNullTime.Valid {
		updateAt = &updateAtNullTime.Time
	}
	var ok bool
	var accrual *big.Float

	if accrualNullString.Valid {
		accrual, ok = new(big.Float).SetString(accrualNullString.String)
	}
	if !ok {
		return &Order{
			login:      login,
			number:     number,
			status:     status,
			accrual:    nil,
			uploadedAt: uploadedAt,
			createdAt:  createdAt,
			updateAt:   updateAt,
		}, nil
	}
	return &Order{
		login:      login,
		number:     number,
		status:     status,
		accrual:    accrual,
		uploadedAt: uploadedAt,
		createdAt:  createdAt,
		updateAt:   updateAt,
	}, nil
}
