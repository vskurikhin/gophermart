/*
 * This file was last modified at 2024-04-20 18:09 by Victor N. Skurikhin.
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

func FuncGetAllOrdersForLogin() func(storage.Storage, string) ([]*Order, error) {
	return func(s storage.Storage, login string) ([]*Order, error) {

		rows, err := s.GetAllForLogin(
			`SELECT o.*, s.status FROM "orders" o JOIN status s ON o.status_id = s.id  WHERE login = $1`,
			login,
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
}

func extractOrder(row pgx.Row) (*Order, error) {

	pLogin, pNumber, status, pAccrual, pUploadedAt, pCreatedAt, pUpdateAt, err := extractOrderTuple(row)

	if err != nil {
		return nil, err
	}

	return &Order{
		login:      *pLogin,
		number:     *pNumber,
		status:     status,
		accrual:    pAccrual,
		uploadedAt: pUploadedAt,
		createdAt:  *pCreatedAt,
		updateAt:   pUpdateAt,
	}, nil
}

func extractOrderTuple(row pgx.Row) (*string, *string, string, *big.Float, *time.Time, *time.Time, *time.Time, error) {

	var statusID int
	var login, number, status string
	var createdAt time.Time
	var accrualNullString sql.NullString
	var uploadedAtNullTime, updateAtNullTime sql.NullTime

	err := row.Scan(
		&login, &number, &statusID, &accrualNullString, &uploadedAtNullTime, &createdAt, &updateAtNullTime, &status,
	)

	if err != nil {
		return nil, nil, status, nil, nil, nil, nil, err
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
		return &login, &number, status, nil, uploadedAt, &createdAt, updateAt, err
	}

	return &login, &number, status, accrual, uploadedAt, &createdAt, updateAt, err
}
