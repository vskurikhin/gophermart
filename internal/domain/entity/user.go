/*
 * This file was last modified at 2024-04-20 17:09 by Victor N. Skurikhin.
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

func NewUser(login string, password *string) *User {
	return &User{login: login, password: password}
}

func (u *User) Login() string {
	return u.login
}

func (u *User) Password() *string {
	return u.password
}

func (u *User) SetPassword(password *string) {
	u.password = password
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) UpdateAt() *time.Time {
	return u.updateAt
}

func (u *User) Insert(s storage.Storage) (*User, error) {
	row, err := s.Save(
		`INSERT INTO "users" (login, password, created_at) VALUES ($1, $2, now()) RETURNING *`,
		u.login, u.password,
	)
	if err != nil {
		return nil, err
	}

	return extractUser(row)
}

func (u *User) AppendInsertTo(a storage.TxArgs) storage.TxArgs {

	t := storage.NewTxArg(
		`INSERT INTO "users" (login, password, created_at) VALUES ($1, $2, now())`,
		u.login, u.password,
	)
	return append(a, t)
}

func (u *User) Save(s storage.Storage) (*User, error) {

	row, err := s.Save(
		`INSERT INTO "users" (login, password, created_at)
             VALUES ($1, $2, now())
             ON CONFLICT (login)
             DO UPDATE SET
               password = $2
             RETURNING *`,
		u.login, u.password,
	)
	if err != nil {
		return nil, err
	}

	return extractUser(row)
}

func FuncGetAllUsers() func(storage.Storage) ([]*User, error) {
	return func(s storage.Storage) ([]*User, error) {
		result := make([]*User, 0)
		return result, nil
	}
}

func FuncGetUser() func(storage.Storage, string) (*User, error) {
	return func(s storage.Storage, login string) (*User, error) {

		row, err := s.GetByLogin(
			`SELECT * FROM "users" WHERE login = $1`,
			login,
		)
		if err != nil {
			return nil, err
		}

		pLogin, pPassword, pCreatedAt, pUpdateAt, err := extractUserTuple(row)

		if err != nil {
			return nil, err
		}

		return &User{
			login:     *pLogin,
			password:  pPassword,
			createdAt: *pCreatedAt,
			updateAt:  pUpdateAt,
		}, nil
	}
}

func extractUser(row pgx.Row) (*User, error) {
	pLogin, pPassword, pCreatedAt, pUpdateAt, err := extractUserTuple(row)

	if err != nil {
		return nil, err
	}

	return &User{
		login:     *pLogin,
		password:  pPassword,
		createdAt: *pCreatedAt,
		updateAt:  pUpdateAt,
	}, nil
}

func extractUserTuple(row pgx.Row) (*string, *string, *time.Time, *time.Time, error) {

	var login string
	var password *string
	var createdAt time.Time
	var passwordNull sql.NullString
	var updateAtNullTime sql.NullTime

	err := row.Scan(&login, &passwordNull, &createdAt, &updateAtNullTime)

	if err != nil {
		return nil, nil, nil, nil, err
	}
	var updateAt *time.Time

	if passwordNull.Valid {
		password = &passwordNull.String
	}
	if updateAtNullTime.Valid {
		updateAt = &updateAtNullTime.Time
	}
	return &login, password, &createdAt, updateAt, err
}
