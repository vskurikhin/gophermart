/*
 * This file was last modified at 2024-04-21 13:25 by Victor N. Skurikhin.
 * mock.go
 * $Id$
 */

package utils

import (
	"github.com/driftprogramming/pgxpoolmock"
	pgx4 "github.com/jackc/pgx/v4"
	pgx5 "github.com/jackc/pgx/v5"
)

func PgxRowsNext(columns []string, values ...interface{}) pgx4.Rows {
	pgxRows := pgxpoolmock.NewRows(columns).
		AddRow(values...).
		ToPgxRows()
	pgxRows.Next()
	return pgxRows
}

type TestRows struct {
}

func (t *TestRows) Scan(dest ...any) error {
	return pgx5.ErrNoRows
}
