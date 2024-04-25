/*
 * This file was last modified at 2024-05-07 14:33 by Victor N. Skurikhin.
 * balance.go
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

type Balance struct {
	login     string
	current   big.Float
	withdrawn big.Float
	createdAt time.Time
	updateAt  *time.Time
}

func GetBalance(s storage.Storage, login string) (*Balance, error) {

	row, err := s.GetByString(
		`SELECT login, current, withdrawn, created_at, update_at
			FROM balance WHERE login = $1`, login,
	)
	if err != nil {
		return nil, err
	}
	return extractBalance(row)
}

func NewBalance(login string, current big.Float) *Balance {
	return &Balance{login: login, current: current}
}

func (b *Balance) Current() big.Float {
	return b.current
}

func (b *Balance) AppendInsertTo(a storage.TxArgs) storage.TxArgs {

	balance, _ := b.current.Float64()
	withdrawn, _ := b.withdrawn.Float64()
	t := storage.NewTxArg(
		`INSERT INTO "balance" (login, current, withdrawn, created_at) VALUES ($1, $2, $3, now())`,
		b.login, balance, withdrawn,
	)
	return append(a, t)
}

func (b *Balance) AppendAccrualTo(a storage.TxArgs, sum *big.Float) storage.TxArgs {

	if sum == nil {
		return a
	}
	current := b.current
	b.current = *b.current.Add(&current, sum)
	accuracy, _ := sum.Float64()
	t := storage.NewTxArg(
		`UPDATE "balance" SET current = current + $1 WHERE login = $2`,
		accuracy, b.login,
	)
	return append(a, t)
}

func (b *Balance) AppendWithdrawTo(a storage.TxArgs, sum *big.Float) storage.TxArgs {

	if sum == nil {
		return a
	}
	current := b.current
	b.current = *b.current.Sub(&current, sum)
	withdraw, _ := sum.Float64()
	t := storage.NewTxArg(
		`UPDATE "balance" SET current = current - $1 WHERE login = $2`,
		withdraw, b.login,
	)
	return append(a, t)
}

func extractBalance(row pgx.Row) (*Balance, error) {

	var login, sCurrent, sWithdrawn string
	var createdAt time.Time
	var updateAtNullTime sql.NullTime

	err := row.Scan(&login, &sCurrent, &sWithdrawn, &createdAt, &updateAtNullTime)

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
	return &Balance{
		login:     login,
		current:   *balance,
		withdrawn: *withdrawn,
		createdAt: createdAt,
		updateAt:  updateAt,
	}, nil
}
