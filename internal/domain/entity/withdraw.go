/*
 * This file was last modified at 2024-04-19 15:58 by Victor N. Skurikhin.
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
	statusID  int
	createdAt time.Time
	updateAt  *time.Time
}

func NewWithdraw(login string, number string, sum big.Float, statusID int) *Withdraw {
	return &Withdraw{login: login, number: number, sum: sum, statusID: statusID}
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

func (w *Withdraw) StatusID() int {
	return w.statusID
}

func (w *Withdraw) SetStatusID(statusID int) {
	w.statusID = statusID
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
		w.login, w.number, sum, w.statusID,
	)
	if err != nil {
		return nil, err
	}

	pLogin, pNumber, pSum, pStatusID, pCreatedAt, pUpdateAt, err := extractWithdraw(row)

	if err != nil {
		return nil, err
	}

	return &Withdraw{
		login:     *pLogin,
		number:    *pNumber,
		sum:       *pSum,
		statusID:  *pStatusID,
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

		pLogin, pNumber, pSum, pStatusID, pCreatedAt, pUpdateAt, err := extractWithdraw(row)

		if err != nil {
			return nil, err
		}

		return &Withdraw{
			login:     *pLogin,
			number:    *pNumber,
			sum:       *pSum,
			statusID:  *pStatusID,
			createdAt: *pCreatedAt,
			updateAt:  pUpdateAt,
		}, nil
	}
}

var zero = big.NewFloat(0.0)

func FuncGetWithdrawSum() func(storage.Storage, string) (*big.Float, error) {
	return func(s storage.Storage, login string) (*big.Float, error) {

		row, err := s.GetByLogin(
			`SELECT sum(sum) FROM withdraw WHERE login = $1`,
			login,
		)
		if err != nil {
			return nil, err
		}

		var sSum string
		err = row.Scan(&sSum)

		if err != nil {
			return zero, nil
		}
		sum, ok := new(big.Float).SetString(sSum)

		if ok {
			return sum, nil
		}
		return zero, nil
	}
}

func extractWithdraw(row pgx.Row) (*string, *string, *big.Float, *int, *time.Time, *time.Time, error) {

	var statusID int
	var login, number, sSum string
	var createdAt time.Time
	var updateAtNullTime sql.NullTime

	err := row.Scan(&login, &number, &sSum, &statusID, &createdAt, &updateAtNullTime)

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
	return &login, &number, sum, &statusID, &createdAt, updateAt, err
}
