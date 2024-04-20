/*
 * This file was last modified at 2024-04-20 17:09 by Victor N. Skurikhin.
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

func NewOrder(login string, number string) *Order {
	return &Order{login: login, number: number}
}

func NewOrderStatus(login string, number string, statusID int) *Order {
	return &Order{login: login, number: number, statusID: statusID}
}

func (o *Order) Login() string {
	return o.login
}

func (o *Order) Number() string {
	return o.number
}

func (o *Order) SetNumber(number string) {
	o.number = number
}

func (o *Order) StatusID() int {
	return o.statusID
}

func (o *Order) SetStatusID(statusID int) {
	o.statusID = statusID
}

func (o *Order) Accrual() *big.Float {
	return o.accrual
}

func (o *Order) SetAccrual(accrual *big.Float) {
	o.accrual = accrual
}

func (o *Order) UploadedAt() *time.Time {
	return o.uploadedAt
}

func (o *Order) CreatedAt() time.Time {
	return o.createdAt
}

func (o *Order) UpdateAt() *time.Time {
	return o.updateAt
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

func (o *Order) Save(s storage.Storage) (*Order, error) {

	var accrualNullFloat64 sql.NullFloat64
	var uploadedAtNullTime sql.NullTime

	if o.accrual != nil {
		accrual, _ := (*o).accrual.Float64()
		accrualNullFloat64.Float64 = accrual
		accrualNullFloat64.Valid = true
	}
	if o.uploadedAt != nil {
		uploadedAtNullTime.Time = *o.uploadedAt
		uploadedAtNullTime.Valid = true
	}
	row, err := s.Save(
		`INSERT INTO "orders"
				    (login, number, status_id, accrual, uploaded_at, created_at)
             VALUES ($1, $2, $3, $4, $5, now())
             ON CONFLICT (login, number)
             DO UPDATE SET
               status_id = $3,
			   accrual = $4,
               uploaded_at = $5
             RETURNING *`,
		o.login, o.number, o.statusID, accrualNullFloat64, uploadedAtNullTime,
	)
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

func FuncGetAllOrders() func(storage.Storage) ([]*Order, error) {
	return func(s storage.Storage) ([]*Order, error) {
		result := make([]*Order, 0)
		return result, nil
	}
}

func FuncGetOrderByNumber() func(storage.Storage, string) (*Order, error) {
	return func(s storage.Storage, number string) (*Order, error) {

		row, err := s.GetByNumber(
			`SELECT * FROM "orders" WHERE number = $1`,
			number,
		)
		if err != nil {
			return nil, err
		}
		return extractOrder(row)
	}
}

func FuncGetOrder() func(storage.Storage, string, string) (*Order, error) {
	return func(s storage.Storage, login, number string) (*Order, error) {

		row, err := s.GetByLoginNumber(
			`SELECT * FROM "orders" WHERE login = $1 AND number = $2`,
			login, number,
		)
		if err != nil {
			return nil, err
		}
		return extractOrder(row)
	}
}

func extractOrder(row pgx.Row) (*Order, error) {

	pLogin, pNumber, pStatusID, pAccrual, pUploadedAt, pCreatedAt, pUpdateAt, err := extractOrderTuple(row)

	if err != nil {
		return nil, err
	}

	return &Order{
		login:      *pLogin,
		number:     *pNumber,
		statusID:   *pStatusID,
		accrual:    pAccrual,
		uploadedAt: pUploadedAt,
		createdAt:  *pCreatedAt,
		updateAt:   pUpdateAt,
	}, nil
}

func extractOrderTuple(row pgx.Row) (*string, *string, *int, *big.Float, *time.Time, *time.Time, *time.Time, error) {

	var statusID int
	var login, number string
	var createdAt time.Time
	var accrualNullString sql.NullString
	var uploadedAtNullTime, updateAtNullTime sql.NullTime

	err := row.Scan(&login, &number, &statusID, &accrualNullString, &uploadedAtNullTime, &createdAt, &updateAtNullTime)

	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
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
		return &login, &number, &statusID, nil, uploadedAt, &createdAt, updateAt, err
	}

	return &login, &number, &statusID, accrual, uploadedAt, &createdAt, updateAt, err
}
