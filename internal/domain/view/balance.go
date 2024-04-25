/*
 * This file was last modified at 2024-04-25 22:06 by Victor N. Skurikhin.
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

type Balance struct {
	login     string
	current   big.Float
	sum       big.Float
	withdrawn big.Float
	createdAt time.Time
	updateAt  *time.Time
}

func GetBalanceWithdraw(s storage.Storage, login string) (*Balance, error) {

	row, err := s.GetByString(
		`SELECT login, current, withdrawn, created_at, update_at, (SELECT sum(sum) AS sum
			FROM withdraw WHERE login = $1) FROM "balance" WHERE login = $1`, login,
	)
	if err != nil {
		return nil, err
	}
	balance, err := extractBalanceWithdrawn(row)

	if utils.IsErrNoRowsInResultSet(err) {
		return &Balance{login: login, current: utils.BigFloatWith0()}, err
	} else if err != nil {
		return nil, err
	}
	return balance, nil
}

func (b *Balance) Current() big.Float {
	return b.current
}

func (b *Balance) Sum() big.Float {
	return b.sum
}

func extractBalanceWithdrawn(row pgx.Row) (*Balance, error) {

	var login, sCurrent, sWithdrawn string
	var createdAt time.Time
	var updateAtNullTime sql.NullTime
	var sumNull sql.NullString
	var sum *big.Float

	err := row.Scan(&login, &sCurrent, &sWithdrawn, &createdAt, &updateAtNullTime, &sumNull)

	if err != nil {
		return nil, err
	}
	balance, ok := new(big.Float).SetString(sCurrent)

	if !ok {
		return nil, errors.New("can't read current")
	}
	withdrawn, ok := new(big.Float).SetString(sWithdrawn)

	if !ok {
		return nil, errors.New("can't read withdrawn")
	}
	var updateAt *time.Time

	if updateAtNullTime.Valid {
		updateAt = &updateAtNullTime.Time
	}
	bigFloatZero := utils.BigFloatWith0()
	if sumNull.Valid {
		sum, ok = new(big.Float).SetString(sumNull.String)
	} else {
		sum = &bigFloatZero
	}
	if !ok {
		return &Balance{
			login:     login,
			current:   *balance,
			sum:       bigFloatZero,
			withdrawn: *withdrawn,
			createdAt: createdAt,
			updateAt:  updateAt,
		}, nil
	}
	return &Balance{
		login:     login,
		current:   *balance,
		sum:       *sum,
		withdrawn: *withdrawn,
		createdAt: createdAt,
		updateAt:  updateAt,
	}, nil
}
