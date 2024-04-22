/*
 * This file was last modified at 2024-04-22 10:44 by Victor N. Skurikhin.
 * consts.go
 * $Id$
 */

package utils

import (
	"database/sql"
	"math/big"
	"time"
)

const (
	StringZero = "0"
)

var (
	bigFloatZero big.Float
	bigIntZero   big.Int
)

func init() {
	f, _ := new(big.Float).SetString(StringZero)
	bigFloatZero = *f
	i := new(big.Int).SetInt64(0)
	bigIntZero = *i
}

func BigFloatWith0() big.Float {
	return bigFloatZero
}

func BigFloatZero() big.Float {
	return big.Float{}
}

func BigIntWith0() big.Int {
	return bigIntZero
}

func SQLNullStringNull() sql.NullString {
	return sql.NullString{}
}

func SQLNullStringZero() sql.NullString {
	return sql.NullString{Valid: true}
}
func SQLNullStringWith0() sql.NullString {
	return sql.NullString{String: StringZero, Valid: true}
}

func SQLNullTimeNull() sql.NullTime {
	return sql.NullTime{}
}

func SQLNullTimeZero() sql.NullTime {
	return sql.NullTime{Valid: true}
}

func TimeEpoch() time.Time {
	return time.Time{}
}

func TimeZero() time.Time {
	return time.Unix(0, 0)
}
