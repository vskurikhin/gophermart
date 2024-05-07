/*
 * This file was last modified at 2024-04-25 22:07 by Victor N. Skurikhin.
 * user.go
 * $Id$
 */

package entity

import (
	"database/sql"
	"github.com/jackc/pgx/v5"
	"github.com/vskurikhin/gophermart/internal/storage"
	"time"
)

type User struct {
	login     string
	password  *string
	createdAt time.Time
	updateAt  *time.Time
}

func GetUser(s storage.Storage, login string) (*User, error) {

	row, err := s.GetByString(
		`SELECT login, password, created_at, update_at
			FROM "users" WHERE login = $1`,
		login,
	)
	if err != nil {
		return nil, err
	}

	return extractUser(row)
}

func NewUser(login string, password *string) *User {
	return &User{login: login, password: password}
}

func (u *User) Login() string {
	return u.login
}

func (u *User) Password() *string {
	return u.password
}

func (u *User) AppendInsertTo(a storage.TxArgs) storage.TxArgs {

	t := storage.NewTxArg(
		`INSERT INTO "users" (login, password, created_at) VALUES ($1, $2, now())`,
		u.login, u.password,
	)
	return append(a, t)
}

func extractUser(row pgx.Row) (*User, error) {

	var login string
	var password *string
	var createdAt time.Time
	var passwordNull sql.NullString
	var updateAtNullTime sql.NullTime

	err := row.Scan(&login, &passwordNull, &createdAt, &updateAtNullTime)

	if err != nil {
		return nil, err
	}
	var updateAt *time.Time

	if passwordNull.Valid {
		password = &passwordNull.String
	}
	if updateAtNullTime.Valid {
		updateAt = &updateAtNullTime.Time
	}
	return &User{
		login:     login,
		password:  password,
		createdAt: createdAt,
		updateAt:  updateAt,
	}, nil
}
