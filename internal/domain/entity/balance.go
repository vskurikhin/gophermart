/*
 * This file was last modified at 2024-04-20 17:09 by Victor N. Skurikhin.
 * balance.go
 * $Id$
 */

package entity

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
	withdrawn big.Float
	createdAt time.Time
	updateAt  *time.Time
}

func NewBalance(login string, current big.Float) *Balance {
	return &Balance{login: login, current: current}
}

func (b *Balance) Login() string {
	return b.login
}

func (b *Balance) Current() big.Float {
	return b.current
}

func (b *Balance) SetCurrent(balance big.Float) {
	b.current = balance
}

func (b *Balance) Withdrawn() big.Float {
	return b.withdrawn
}

func (b *Balance) SetWithdrawn(withdrawn big.Float) {
	b.withdrawn = withdrawn
}

func (b *Balance) CreatedAt() time.Time {
	return b.createdAt
}

func (b *Balance) UpdateAt() *time.Time {
	return b.updateAt
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

	accuracy, _ := sum.Float64()
	t := storage.NewTxArg(
		`UPDATE "balance" SET current = current + $1 WHERE login = $2`,
		accuracy, b.login,
	)
	return append(a, t)
}

func (b *Balance) AppendWithdrawTo(a storage.TxArgs, sum *big.Float) storage.TxArgs {

	withdraw, _ := sum.Float64()
	t := storage.NewTxArg(
		`UPDATE "balance" SET current = current - $1 WHERE login = $2`,
		withdraw, b.login,
	)
	return append(a, t)
}

func (b *Balance) Save(s storage.Storage) (*Balance, error) {

	balance, _ := b.current.Float64()
	withdrawn, _ := b.withdrawn.Float64()
	row, err := s.Save(
		`INSERT INTO "balance"
				    (login, current, withdrawn, created_at)
             VALUES ($1, $2, $3, now())
             ON CONFLICT (login)
             DO UPDATE SET
               current = $2,
               withdrawn = $3
             RETURNING *`,
		b.login, balance, withdrawn,
	)
	if err != nil {
		return nil, err
	}

	pLogin, pBalance, pWithdrawn, pCreatedAt, pUpdateAt, err := extractBalance(row)

	if err != nil {
		return nil, err
	}

	return &Balance{
		login:     *pLogin,
		current:   *pBalance,
		withdrawn: *pWithdrawn,
		createdAt: *pCreatedAt,
		updateAt:  pUpdateAt,
	}, nil
}

func FuncGetAllBalances() func(storage.Storage) ([]*Balance, error) {
	return func(s storage.Storage) ([]*Balance, error) {
		result := make([]*Balance, 0)
		return result, nil
	}
}

func FuncGetBalanceWithdraw() func(storage.Storage, string) (*Balance, *big.Float, error) {
	return func(s storage.Storage, login string) (*Balance, *big.Float, error) {

		row, err := s.GetByLogin(
			`SELECT *, (SELECT sum(sum) FROM withdraw WHERE login = $1) FROM "balance" WHERE login = $1`,
			login,
		)

		if err != nil {
			return nil, zero, err
		}

		_, pCurrent, pWithdrawn, pCreatedAt, pUpdateAt, sum, err := extractBalanceWithdrawn(row)

		if utils.IsErrNoRowsInResultSet(err) {
			return &Balance{login: login, current: *zero}, zero, err
		} else if err != nil {
			return nil, nil, err
		}

		return &Balance{
			login:     login,
			current:   *pCurrent,
			withdrawn: *pWithdrawn,
			createdAt: *pCreatedAt,
			updateAt:  pUpdateAt,
		}, sum, nil
	}
}

func FuncGetBalance() func(storage.Storage, string) (*Balance, error) {
	return func(s storage.Storage, login string) (*Balance, error) {

		row, err := s.GetByLogin("SELECT * FROM balance WHERE login = $1", login)

		if err != nil {
			return nil, err
		}

		_, pCurrent, pWithdrawn, pCreatedAt, pUpdateAt, err := extractBalance(row)

		if err != nil {
			return nil, err
		}

		return &Balance{
			login:     login,
			current:   *pCurrent,
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
	}

	if !ok {
		return &login, balance, withdrawn, &createdAt, updateAt, zero, nil
	}
	return &login, balance, withdrawn, &createdAt, updateAt, sum, err
}

func extractBalance(row pgx.Row) (*string, *big.Float, *big.Float, *time.Time, *time.Time, error) {

	var login, sCurrent, sWithdrawn string
	var createdAt time.Time
	var updateAtNullTime sql.NullTime

	err := row.Scan(&login, &sCurrent, &sWithdrawn, &createdAt, &updateAtNullTime)

	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	balance, ok := new(big.Float).SetString(sCurrent)

	if !ok {
		return nil, nil, nil, nil, nil, errors.New("can't read current")
	}
	withdrawn, ok := new(big.Float).SetString(sWithdrawn)

	if !ok {
		return nil, nil, nil, nil, nil, errors.New("can't read withdrawn")
	}
	var updateAt *time.Time

	if updateAtNullTime.Valid {
		updateAt = &updateAtNullTime.Time
	}
	return &login, balance, withdrawn, &createdAt, updateAt, err
}
