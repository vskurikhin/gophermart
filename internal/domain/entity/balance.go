/*
 * This file was last modified at 2024-04-12 13:23 by Victor N. Skurikhin.
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
	balance   big.Float
	withdrawn big.Float
	createdAt time.Time
	updateAt  *time.Time
}

func NewBalance(login string, balance big.Float) *Balance {
	return &Balance{login: login, balance: balance}
}

func NewBalanceWithdrawn(login string, balance big.Float, withdrawn big.Float) *Balance {
	return &Balance{login: login, balance: balance, withdrawn: withdrawn}
}

func (b *Balance) Login() string {
	return b.login
}

func (b *Balance) Balance() big.Float {
	return b.balance
}

func (b *Balance) SetBalance(balance big.Float) {
	b.balance = balance
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

func (b *Balance) Save(s storage.Storage) (*Balance, error) {

	balance, _ := b.balance.Float64()
	withdrawn, _ := b.withdrawn.Float64()
	row, err := s.Save(
		`INSERT INTO balance
				    (login, balance, withdrawn, created_at)
             VALUES ($1, $2, $3, now())
             ON CONFLICT (login)
             DO UPDATE SET
               balance = $2,
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
		balance:   *pBalance,
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

func FuncGetBalance() func(storage.Storage, string) (*Balance, error) {
	return func(s storage.Storage, login string) (*Balance, error) {

		row, err := s.GetByLogin("SELECT * FROM balance WHERE login = $1", login)

		if err != nil {
			return nil, err
		}

		_, pBalance, pWithdrawn, pCreatedAt, pUpdateAt, err := extractBalance(row)

		if err != nil {
			return nil, err
		}

		return &Balance{
			login:     login,
			balance:   *pBalance,
			withdrawn: *pWithdrawn,
			createdAt: *pCreatedAt,
			updateAt:  pUpdateAt,
		}, nil
	}
}

func extractBalance(row pgx.Row) (*string, *big.Float, *big.Float, *time.Time, *time.Time, error) {

	var login, sBalance, sWithdrawn string
	var createdAt time.Time
	var updateAtNullTime sql.NullTime

	err := row.Scan(&login, &sBalance, &sWithdrawn, &createdAt, &updateAtNullTime)

	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	balance, ok := new(big.Float).SetString(sBalance)

	if !ok {
		return nil, nil, nil, nil, nil, errors.New("can't read balance")
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