/*
 * This file was last modified at 2024-04-19 19:59 by Victor N. Skurikhin.
 * tx_args.go
 * $Id$
 */

package storage

type TxArg struct {
	sql    string
	values []interface{}
}

func NewTxArg(sql string, values ...interface{}) *TxArg {
	return &TxArg{sql: sql, values: values}
}

type TxArgs []*TxArg

func (t *TxArgs) Append(sql string, values ...interface{}) {
	txArg := NewTxArg(sql, values...)
	*t = append(*t, txArg)
}
