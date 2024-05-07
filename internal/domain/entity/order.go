/*
 * This file was last modified at 2024-04-25 21:52 by Victor N. Skurikhin.
 * order.go
 * $Id$
 */

package entity

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
	statusID   int
	accrual    *big.Float
	uploadedAt *time.Time
	createdAt  time.Time
	updateAt   *time.Time
}

func GetOrderByNumber(s storage.Storage, number string) (*Order, error) {

	row, err := s.GetByString(
		`SELECT login, number, status_id, accrual, uploaded_at, created_at, update_at
			FROM "orders" WHERE number = $1`,
		number,
	)
	if err != nil {
		return nil, err
	}
	return extractOrder(row)
}

func NewOrder(login string, number string) *Order {
	return &Order{login: login, number: number}
}

func (o *Order) Login() string {
	return o.login
}

func (o *Order) Accrual() *big.Float {
	return o.accrual
}

func (o *Order) SetAccrual(accrual *big.Float) {
	o.accrual = accrual
}

func (o *Order) Insert(s storage.Storage) (*Order, error) {

	var err error
	var row pgx.Row

	if o.statusID > 0 {
		row, err = s.Save(
			`INSERT INTO "orders"
				    (login, number, status_id, uploaded_at, created_at)
             VALUES ($1, $2, $3, now(), now())
             RETURNING *`,
			o.login, o.number, o.statusID,
		)
	} else {
		row, err = s.Save(
			`INSERT INTO "orders"
				    (login, number, status_id, uploaded_at, created_at)
             VALUES ($1, $2, (SELECT id FROM status WHERE status = 'NEW'), now(), now())
             RETURNING *`,
			o.login, o.number,
		)
	}
	if err != nil {
		return nil, err
	}
	return extractOrder(row)
}

func (o *Order) AppendAccrualTo(a storage.TxArgs) storage.TxArgs {

	var accuracy sql.NullFloat64
	if o.accrual != nil {
		b := o.accrual
		a, _ := b.Float64()
		accuracy.Float64 = a
		accuracy.Valid = true
	}
	t := storage.NewTxArg(
		`UPDATE "orders" SET accrual = $1 WHERE number = $2`,
		accuracy, o.number,
	)
	return append(a, t)
}

func (o *Order) AppendSetStatusTo(a storage.TxArgs, status string) storage.TxArgs {

	t := storage.NewTxArg(
		`UPDATE "orders" SET status_id = (SELECT id FROM status WHERE status = $1) WHERE number = $2`,
		status, o.number,
	)
	return append(a, t)
}

func extractOrder(row pgx.Row) (*Order, error) {

	var statusID int
	var login, number string
	var createdAt time.Time
	var accrualNullString sql.NullString
	var uploadedAtNullTime, updateAtNullTime sql.NullTime

	err := row.Scan(&login, &number, &statusID, &accrualNullString, &uploadedAtNullTime, &createdAt, &updateAtNullTime)

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
			statusID:   statusID,
			accrual:    nil,
			uploadedAt: uploadedAt,
			createdAt:  createdAt,
			updateAt:   updateAt,
		}, nil
	}
	return &Order{
		login:      login,
		number:     number,
		statusID:   statusID,
		accrual:    accrual,
		uploadedAt: uploadedAt,
		createdAt:  createdAt,
		updateAt:   updateAt,
	}, nil
}
