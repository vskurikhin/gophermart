/*
 * This file was last modified at 2024-04-13 17:14 by Victor N. Skurikhin.
 * withdraw.go
 * $Id$
 */

package entity

import (
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/vskurikhin/gophermart/internal/storage"
	"math/big"
	"time"
)

type Withdraw struct {
	login     string
	number    string
	sum       big.Float
	statusId  int
	createdAt time.Time
	updateAt  *time.Time
}

func NewWithdraw(login string, number string, sum big.Float, statusId int) *Withdraw {
	return &Withdraw{login: login, number: number, sum: sum, statusId: statusId}
}

func NewWithdrawWithTime(login string, number string, sum big.Float, statusId int, createdAt time.Time, updateAt *time.Time) *Withdraw {
	return &Withdraw{login: login, number: number, sum: sum, statusId: statusId, createdAt: createdAt, updateAt: updateAt}
}

func (w *Withdraw) Login() string {
	return w.login
}

func (w *Withdraw) Number() string {
	return w.number
}

func (w *Withdraw) Sum() big.Float {
	return w.sum
}

func (w *Withdraw) SetSum(sum big.Float) {
	w.sum = sum
}

func (w *Withdraw) StatusId() int {
	return w.statusId
}

func (w *Withdraw) SetStatusId(statusId int) {
	w.statusId = statusId
}

func (w *Withdraw) CreatedAt() time.Time {
	return w.createdAt
}

func (w *Withdraw) UpdateAt() *time.Time {
	return w.updateAt
}

func (w *Withdraw) Save(s storage.Storage) (*Withdraw, error) {

	sum, _ := w.sum.Float64()
	row, err := s.Save(
		`INSERT INTO withdraw
				    (login, number, sum, status_id, created_at)
             VALUES ($1, $2, $3, $4, now())
             ON CONFLICT (login, number)
             DO UPDATE SET
               sum = $3,
               status_id = $4
		     RETURNING *`,
		w.login, w.number, sum, w.statusId,
	)
	if err != nil {
		return nil, err
	}

	pLogin, pNumber, pSum, pStatusId, pCreatedAt, pUpdateAt, err := extractWithdraw(row)

	if err != nil {
		return nil, err
	}

	return &Withdraw{
		login:     *pLogin,
		number:    *pNumber,
		sum:       *pSum,
		statusId:  *pStatusId,
		createdAt: *pCreatedAt,
		updateAt:  pUpdateAt,
	}, nil
}

func FuncGetAllWithdraw() func(storage.Storage) ([]*Withdraw, error) {
	return func(s storage.Storage) ([]*Withdraw, error) {
		result := make([]*Withdraw, 0)
		return result, nil
	}
}

func FuncGetWithdraw() func(storage.Storage, string, string) (*Withdraw, error) {
	return func(s storage.Storage, login, number string) (*Withdraw, error) {

		row, err := s.GetByLoginNumber(
			`SELECT * FROM withdraw WHERE login = $1 AND number = $2`,
			login, number,
		)
		if err != nil {
			return nil, err
		}

		pLogin, pNumber, pSum, pStatusId, pCreatedAt, pUpdateAt, err := extractWithdraw(row)

		if err != nil {
			return nil, err
		}

		return &Withdraw{
			login:     *pLogin,
			number:    *pNumber,
			sum:       *pSum,
			statusId:  *pStatusId,
			createdAt: *pCreatedAt,
			updateAt:  pUpdateAt,
		}, nil
	}
}

func extractWithdraw(row pgx.Row) (*string, *string, *big.Float, *int, *time.Time, *time.Time, error) {

	var statusId int
	var login, number, sSum string
	var createdAt time.Time
	var updateAtNullTime sql.NullTime

	err := row.Scan(&login, &number, &sSum, &statusId, &createdAt, &updateAtNullTime)

	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}
	sum, ok := new(big.Float).SetString(sSum)

	if !ok {
		return nil, nil, nil, nil, nil, nil, errors.New("can't read sum")
	}
	var updateAt *time.Time

	if updateAtNullTime.Valid {
		updateAt = &updateAtNullTime.Time
	}
	return &login, &number, sum, &statusId, &createdAt, updateAt, err
}
