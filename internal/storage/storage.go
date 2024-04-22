/*
 * This file was last modified at 2024-04-21 00:49 by Victor N. Skurikhin.
 * storage.go
 * $Id$
 */

package storage

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type Storage interface {
	GetAllForString(sql, login string) (pgx.Rows, error)

	GetByID(sql string, id int) (pgx.Row, error)

	GetByString(sql, str string) (pgx.Row, error)

	GetByStr1Str2(sql, str1, str2 string) (pgx.Row, error)

	Save(sql string, values ...any) (pgx.Row, error)

	Transaction(args TxArgs) error

	WithContext(ctx context.Context) Storage
}
