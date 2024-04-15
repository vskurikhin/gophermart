/*
 * This file was last modified at 2024-04-15 16:57 by Victor N. Skurikhin.
 * storage.go
 * $Id$
 */

package storage

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type Storage interface {
	GetAll(sql string) (pgx.Rows, error)

	GetByID(sql string, id int) (pgx.Row, error)

	GetByLogin(sql, login string) (pgx.Row, error)

	GetByLoginNumber(sql, login, number string) (pgx.Row, error)

	Save(sql string, values ...any) (pgx.Row, error)

	WithContext(ctx context.Context) Storage
}
