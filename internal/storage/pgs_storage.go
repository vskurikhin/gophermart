/*
 * This file was last modified at 2024-04-21 00:49 by Victor N. Skurikhin.
 * pgs_storage.go
 * $Id$
 */

package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vskurikhin/gophermart/internal/logger"
	"github.com/vskurikhin/gophermart/internal/utils"
	"go.uber.org/zap"
	"time"
)

const (
	increase = 2
	timeout  = 10
	tries    = 3
)

type PgsStorage struct {
	ctx  context.Context
	log  *zap.Logger
	pool *pgxpool.Pool
}

func NewPgsStorage() Storage {
	if pool, ok := GetDB().DBPool(); ok {
		return &PgsStorage{ctx: context.Background(), log: logger.Get(), pool: pool}
	}
	panic(fmt.Errorf("can't create DB pool"))
}

func (p *PgsStorage) WithContext(ctx context.Context) Storage {
	return &PgsStorage{ctx: ctx, log: p.log, pool: p.pool}
}

func (p *PgsStorage) GetAllForString(sql, str string) (pgx.Rows, error) {
	const funcName = "PgsStorage.GetAllForString"
	defer utils.TraceInOut(p.ctx, funcName, "%s, login", sql, str)()
	return p.sqlRows(funcName, sql, str)
}

func (p *PgsStorage) GetByID(sql string, id int) (pgx.Row, error) {
	const funcName = "PgsStorage.GetByID"
	defer utils.TraceInOut(p.ctx, funcName, "%s, %d", sql, id)()
	return p.sqlRow(funcName, sql, id)
}

func (p *PgsStorage) GetByString(sql, str string) (pgx.Row, error) {
	const funcName = "PgsStorage.GetByString"
	defer utils.TraceInOut(p.ctx, funcName, "%s, %s", sql, str)()
	return p.sqlRow(funcName, sql, str)
}

func (p *PgsStorage) GetByStr1Str2(sql, str1, str2 string) (pgx.Row, error) {
	const funcName = "PgsStorage.GetByStr1Str2"
	defer utils.TraceInOut(p.ctx, funcName, "%s, %s, %s", sql, str1, str2)()
	return p.sqlRow(funcName, sql, str1, str2)
}

func (p *PgsStorage) Save(sql string, values ...any) (pgx.Row, error) {
	const funcName = "PgsStorage.Save"
	defer utils.TraceInOut(p.ctx, funcName, "%s, %s", sql, values)()
	return p.sqlRow(funcName, sql, values...)
}

func (p *PgsStorage) Transaction(args TxArgs) error {
	const funcName = "PgsStorage.Transaction"
	defer utils.TraceInOut(p.ctx, funcName, "%v", args)()
	return p.transactionConnectionRow(funcName, args)
}

func (p *PgsStorage) sqlRow(name, sql string, values ...any) (pgx.Row, error) {

	defer func() {
		if r := recover(); r != nil {
			p.log.Error(name, utils.LogCtxRecoverFields(p.ctx, r)...)
		}
	}()

	ctx, cancel := context.WithTimeout(p.ctx, time.Duration(timeout)*time.Second)
	defer func() {
		cancel()
		ctx.Done()
	}()

	conn, err := p.pool.Acquire(ctx)

	for i := 1; err != nil && i < tries*increase; i += increase {
		time.Sleep(time.Duration(i) * time.Second)
		p.log.Warn(name, utils.LogCtxReasonErrFields(ctx, "retry pool acquire", err)...)
		conn, err = p.pool.Acquire(ctx)
	}
	defer func() {
		if conn != nil {
			conn.Release()
		}
	}()

	if conn == nil || err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	return conn.QueryRow(ctx, sql, values...), nil
}

func (p *PgsStorage) sqlRows(name, sql string, values ...any) (pgx.Rows, error) {

	defer func() {
		if r := recover(); r != nil {
			p.log.Error(name, utils.LogCtxRecoverFields(p.ctx, r)...)
		}
	}()

	ctx, cancel := context.WithTimeout(p.ctx, time.Duration(timeout)*time.Second)
	defer func() {
		cancel()
		ctx.Done()
	}()

	conn, err := p.pool.Acquire(ctx)

	for i := 1; err != nil && i < tries*increase; i += increase {
		time.Sleep(time.Duration(i) * time.Second)
		p.log.Warn(name, utils.LogCtxReasonErrFields(ctx, "retry pool acquire", err)...)
		conn, err = p.pool.Acquire(ctx)
	}
	defer func() {
		if conn != nil {
			conn.Release()
		}
	}()

	if conn == nil || err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	return conn.Query(ctx, sql, values...)
}

func (p *PgsStorage) transactionConnectionRow(name string, args TxArgs) error {

	defer func() {
		if r := recover(); r != nil {
			p.log.Error(name, utils.LogCtxRecoverFields(p.ctx, r)...)
		}
	}()

	ctx, cancel := context.WithTimeout(p.ctx, time.Duration(timeout)*time.Second)
	defer func() {
		cancel()
		ctx.Done()
	}()

	conn, err := p.pool.Acquire(ctx)

	for i := 1; err != nil && i < tries*increase; i += increase {
		time.Sleep(time.Duration(i) * time.Second)
		p.log.Warn(name, utils.LogCtxReasonErrFields(ctx, "retry pool acquire", err)...)
		conn, err = p.pool.Acquire(ctx)
	}
	defer func() {
		if conn != nil {
			conn.Release()
		}
	}()

	if conn == nil || err != nil {
		return fmt.Errorf("%v", err)
	}
	return p.transactionExecs(ctx, conn, args)
}

func (p *PgsStorage) transactionExecs(ctx context.Context, conn *pgxpool.Conn, args TxArgs) error {

	tx, err := conn.Begin(ctx)
	//goland:noinspection GoUnhandledErrorResult
	defer tx.Rollback(ctx)

	if err != nil {
		return err
	}
	for _, arg := range args {
		ct, err := tx.Exec(ctx, arg.sql, arg.values...)
		if err != nil {
			return err
		}
		p.log.Debug("transactionExecs", zap.Int64("RowsAffected", ct.RowsAffected()))
	}
	return tx.Commit(ctx)
}
