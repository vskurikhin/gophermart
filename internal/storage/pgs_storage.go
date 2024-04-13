/*
 * This file was last modified at 2024-04-13 19:09 by Victor N. Skurikhin.
 * pgs_storage.go
 * $Id$
 */

package storage

import (
	"context"
	"errors"
	"fmt"
	uuid4 "github.com/google/uuid"
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
	pool *pgxpool.Pool
	log  *zap.Logger
}

func NewPgsStorage(pool *pgxpool.Pool) *PgsStorage {
	return &PgsStorage{ctx: context.Background(), pool: pool, log: logger.Get()}
}

func (p *PgsStorage) WithContext(ctx context.Context) Storage {
	return &PgsStorage{ctx: ctx, pool: p.pool, log: p.log}
}

func (p *PgsStorage) GetAll(sql string) (pgx.Rows, error) {
	const funcName = "PgsStorage.GetAll"
	defer utils.TraceInOut(p.ctx, funcName, "%s", sql)()
	return nil, errors.New("not implemented")
}

func (p *PgsStorage) GetById(sql string, id int) (pgx.Row, error) {
	const funcName = "PgsStorage.GetById"
	defer utils.TraceInOut(p.ctx, funcName, "%s, %d", sql, id)()
	return p.sqlRow(funcName, sql, id)
}

func (p *PgsStorage) GetByLogin(sql, login string) (pgx.Row, error) {
	const funcName = "PgsStorage.GetByLogin"
	defer utils.TraceInOut(p.ctx, funcName, "%s, %s", sql, login)()
	return p.sqlRow(funcName, sql, login)
}

func (p *PgsStorage) GetByLoginNumber(sql, login, number string) (pgx.Row, error) {
	const funcName = "PgsStorage.GetByLoginNumber"
	defer utils.TraceInOut(p.ctx, funcName, "%s, %s, %s", sql, login, number)()
	return p.sqlRow(funcName, sql, login, number)
}

func (p *PgsStorage) Save(sql string, values ...any) (pgx.Row, error) {
	const funcName = "PgsStorage.Save"
	defer utils.TraceInOut(p.ctx, funcName, "%s, %s", sql, values)()
	return p.sqlRow(funcName, sql, values...)
}

func (p *PgsStorage) sqlRow(name, sql string, values ...any) (pgx.Row, error) {

	defer func() {
		if r := recover(); r != nil {
			var uuid uuid4.UUID
			if id, ok := p.ctx.Value("uuid").(uuid4.UUID); ok {
				uuid = id
			}
			p.log.Warn(fmt.Sprintf(
				"[%s] %s, error: %v",
				uuid, name, r,
			))
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
		if uuid, ok := ctx.Value("uuid").(uuid4.UUID); ok {
			p.log.Warn(fmt.Sprintf(
				"[%s] %s, retry pool acquire error: %v, time: %v",
				uuid, name, err, time.Now(),
			))
		}
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
