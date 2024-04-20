/*
 * This file was last modified at 2024-04-20 17:54 by Victor N. Skurikhin.
 * balance.go
 * $Id$
 */

package view

import (
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/vskurikhin/gophermart/internal/storage"
	"github.com/vskurikhin/gophermart/internal/utils"
	"math/big"
	"time"
)

var zero = big.NewFloat(0.0)

type Balance struct {
	login     string
	current   big.Float
	sum       big.Float
	withdrawn big.Float
	createdAt time.Time
	updateAt  *time.Time
}

func (b *Balance) Current() big.Float {
	return b.current
}

func (b *Balance) Sum() big.Float {
	return b.sum
}

func FuncGetBalanceWithdraw() func(storage.Storage, string) (*Balance, error) {
	return func(s storage.Storage, login string) (*Balance, error) {

		row, err := s.GetByLogin(
			`SELECT *, (SELECT sum(sum) FROM withdraw WHERE login = $1) FROM "balance" WHERE login = $1`,
			login,
		)

		if err != nil {
			return nil, err
		}

		_, pCurrent, pWithdrawn, pCreatedAt, pUpdateAt, pSum, err := extractBalanceWithdrawn(row)

		if utils.IsErrNoRowsInResultSet(err) {
			return &Balance{login: login, current: *zero}, err
		} else if err != nil {
			return nil, err
		}

		return &Balance{
			login:     login,
			current:   *pCurrent,
			sum:       *pSum,
			withdrawn: *pWithdrawn,
			createdAt: *pCreatedAt,
			updateAt:  pUpdateAt,
		}, nil
	}
}

func extractBalanceWithdrawn(row pgx.Row) (*string, *big.Float, *big.Float, *time.Time, *time.Time, *big.Float, error) {

	var login, sCurrent, sWithdrawn string
	var createdAt time.Time
	var updateAtNullTime sql.NullTime
	var sumNull sql.NullString
	var sum *big.Float

	err := row.Scan(&login, &sCurrent, &sWithdrawn, &createdAt, &updateAtNullTime, &sumNull)

	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}
	balance, ok := new(big.Float).SetString(sCurrent)

	if !ok {
		return nil, nil, nil, nil, nil, nil, errors.New("can't read current")
	}
	withdrawn, ok := new(big.Float).SetString(sWithdrawn)

	if !ok {
		return nil, nil, nil, nil, nil, nil, errors.New("can't read withdrawn")
	}
	var updateAt *time.Time

	if updateAtNullTime.Valid {
		updateAt = &updateAtNullTime.Time
	}
	if sumNull.Valid {
		sum, ok = new(big.Float).SetString(sumNull.String)
	} else {
		sum = zero
	}
	if !ok {
		return &login, balance, withdrawn, &createdAt, updateAt, zero, nil
	}
	return &login, balance, withdrawn, &createdAt, updateAt, sum, err
}
