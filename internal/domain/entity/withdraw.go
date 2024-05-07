/*
 * This file was last modified at 2024-04-25 21:59 by Victor N. Skurikhin.
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
	login       string
	number      string
	sum         big.Float
	statusID    int
	processedAt *time.Time
	createdAt   time.Time
	updateAt    *time.Time
}

func GetAllWithdrawalsByLogin(s storage.Storage, login string) ([]*Withdraw, error) {
	rows, err := s.GetAllForString(
		`SELECT login, number, sum, status_id, processed_at, created_at, update_at
		FROM withdraw WHERE login = $1 ORDER BY processed_at, created_at`, login,
	)
	if err != nil {
		return nil, err
	}
	result := make([]*Withdraw, 0)

	for rows.Next() {
		order, err := extractWithdraw(rows)
		if err != nil {
			return result, err
		}
		result = append(result, order)
	}
	return result, nil
}

func NewWithdraw(login string, number string, sum big.Float) *Withdraw {
	return &Withdraw{login: login, number: number, sum: sum}
}

func (w *Withdraw) Number() string {
	return w.number
}

func (w *Withdraw) Sum() big.Float {
	return w.sum
}

func (w *Withdraw) ProcessedAt() *time.Time {
	return w.processedAt
}

func (w *Withdraw) CreatedAt() time.Time {
	return w.createdAt
}

func (w *Withdraw) AppendInsertTo(a storage.TxArgs) storage.TxArgs {

	sum, _ := w.sum.Float64()
	t := storage.NewTxArg(
		`INSERT INTO withdraw (login, number, sum, status_id, created_at)
 		VALUES ($1, $2, $3, (SELECT id FROM status WHERE status = 'NEW'), now())`,
		w.login, w.number, sum,
	)
	return append(a, t)
}

func extractWithdraw(row pgx.Row) (*Withdraw, error) {

	var statusID int
	var login, number, sSum string
	var createdAt time.Time
	var processedAtNullTime, updateAtNullTime sql.NullTime

	err := row.Scan(&login, &number, &sSum, &statusID, &processedAtNullTime, &createdAt, &updateAtNullTime)

	if err != nil {
		return nil, err
	}
	sum, ok := new(big.Float).SetString(sSum)

	if !ok {
		return nil, errors.New("can't read sum")
	}
	var updateAt *time.Time

	if updateAtNullTime.Valid {
		updateAt = &updateAtNullTime.Time
	}
	var processedAt *time.Time

	if processedAtNullTime.Valid {
		processedAt = &processedAtNullTime.Time
	}
	return &Withdraw{
		login:       login,
		number:      number,
		sum:         *sum,
		statusID:    statusID,
		processedAt: processedAt,
		createdAt:   createdAt,
		updateAt:    updateAt,
	}, nil
}
