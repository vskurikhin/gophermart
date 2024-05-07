/*
 * This file was last modified at 2024-04-25 21:15 by Victor N. Skurikhin.
 * storage.go
 * $Id$
 */

package storage

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type Row interface {
	Scan(dest ...any) error
}

type Rows interface {
	pgx.Rows
}

type Storage interface {
	GetAllForString(sql, login string) (Rows, error)

	GetByID(sql string, id int) (Row, error)

	GetByString(sql, str string) (Row, error)

	GetByStr1Str2(sql, str1, str2 string) (Row, error)

	Save(sql string, values ...any) (Row, error)

	Transaction(args TxArgs) error

	WithContext(ctx context.Context) Storage
}
