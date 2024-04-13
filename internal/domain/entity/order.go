/*
 * This file was last modified at 2024-04-13 17:14 by Victor N. Skurikhin.
 * order.go
 * $Id$
 */

package entity

import (
	"database/sql"
	"github.com/jackc/pgx/v5"
	"github.com/vskurikhin/gophermart/internal/storage"
	"time"
)

type Order struct {
	login      string
	number     string
	statusId   int
	uploadedAt *time.Time
	createdAt  time.Time
	updateAt   *time.Time
}

func NewOrder(login string, number string, statusId int) *Order {
	return &Order{login: login, number: number, statusId: statusId}
}

func NewUploadedOrder(login string, number string, statusId int, uploadedAt *time.Time) *Order {
	return &Order{login: login, number: number, statusId: statusId, uploadedAt: uploadedAt}
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

func (o *Order) StatusId() int {
	return o.statusId
}

func (o *Order) SetStatusId(statusId int) {
	o.statusId = statusId
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

func (o *Order) Save(s storage.Storage) (*Order, error) {

	var uploadedAtNullTime sql.NullTime

	if o.uploadedAt != nil {
		uploadedAtNullTime.Time = *o.uploadedAt
		uploadedAtNullTime.Valid = true
	}

	row, err := s.Save(
		`INSERT INTO "order"
				    (login, number, status_id, uploaded_at, created_at)
             VALUES ($1, $2, $3, $4, now())
             ON CONFLICT (login, number)
             DO UPDATE SET
               status_id = $3,
               uploaded_at = $4
             RETURNING *`,
		o.login, o.number, o.statusId, uploadedAtNullTime,
	)
	if err != nil {
		return nil, err
	}

	pLogin, pNumber, pStatusId, pUploadedAt, pCreatedAt, pUpdateAt, err := extractOrder(row)

	if err != nil {
		return nil, err
	}

	return &Order{
		login:      *pLogin,
		number:     *pNumber,
		statusId:   *pStatusId,
		uploadedAt: pUploadedAt,
		createdAt:  *pCreatedAt,
		updateAt:   pUpdateAt,
	}, nil
}

func FuncGetAllOrders() func(storage.Storage) ([]*Order, error) {
	return func(s storage.Storage) ([]*Order, error) {
		result := make([]*Order, 0)
		return result, nil
	}
}

func FuncGetOrder() func(storage.Storage, string, string) (*Order, error) {
	return func(s storage.Storage, login, number string) (*Order, error) {

		row, err := s.GetByLoginNumber(
			`SELECT * FROM "order" WHERE login = $1 AND number = $2`,
			login, number,
		)
		if err != nil {
			return nil, err
		}

		pLogin, pNumber, pStatusId, pUploadedAt, pCreatedAt, pUpdateAt, err := extractOrder(row)

		if err != nil {
			return nil, err
		}

		return &Order{
			login:      *pLogin,
			number:     *pNumber,
			statusId:   *pStatusId,
			uploadedAt: pUploadedAt,
			createdAt:  *pCreatedAt,
			updateAt:   pUpdateAt,
		}, nil
	}
}

func extractOrder(row pgx.Row) (*string, *string, *int, *time.Time, *time.Time, *time.Time, error) {

	var statusId int
	var login, number string
	var createdAt time.Time
	var uploadedAtNullTime, updateAtNullTime sql.NullTime

	err := row.Scan(&login, &number, &statusId, &uploadedAtNullTime, &createdAt, &updateAtNullTime)

	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}
	var uploadedAt *time.Time

	if uploadedAtNullTime.Valid {
		uploadedAt = &uploadedAtNullTime.Time
	}
	var updateAt *time.Time

	if updateAtNullTime.Valid {
		updateAt = &updateAtNullTime.Time
	}
	return &login, &number, &statusId, uploadedAt, &createdAt, updateAt, err
}
