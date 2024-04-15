/*
 * This file was last modified at 2024-04-15 21:27 by Victor N. Skurikhin.
 * db_pool.go
 * $Id$
 */

package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vskurikhin/gophermart/internal/env"
	"github.com/vskurikhin/gophermart/internal/logger"
	"go.uber.org/zap"
	"sync"
	"time"
)

var dbPoolHealthInstance DB
var once = new(sync.Once)

type DB interface {
	DBPool() (*pgxpool.Pool, bool)
	dbPing() error
}

type dbPoolHealth struct {
	sync.RWMutex
	log    *zap.Logger
	pool   *pgxpool.Pool
	status bool
}

func GetDB() DB {
	once.Do(func() {
		cfg := env.GetConfig()
		d := new(dbPoolHealth)
		d.log = logger.Get()
		d.pool = newPgxPool(cfg.DataBaseDSN())
		d.status = true
		go d.checkStatus()
		dbPoolHealthInstance = d
	})
	return dbPoolHealthInstance
}

func newPgxPool(dataBaseDSN string) *pgxpool.Pool {

	log := logger.Get()
	config, err := pgxpool.ParseConfig(dataBaseDSN)

	if err != nil {
		panic(err)
	}
	log.Debug("DBConnect config parsed")

	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		log.Debug("Acquire connect ping...")
		if err = conn.Ping(ctx); err != nil {
			panic(err)
		}
		log.Debug("Acquire connect Ok")
		return nil
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)

	if err != nil {
		panic(err)
	}
	log.Debug("NewWithConfig pool created")
	_, err = pool.Acquire(context.Background())

	if err != nil {
		panic(err)
	}
	log.Debug("Acquire pool Ok")

	return pool
}

func (h *dbPoolHealth) DBPool() (*pgxpool.Pool, bool) {
	h.RLock()
	defer h.RUnlock()
	return h.pool, h.status
}

func (h *dbPoolHealth) dbPing() error {
	h.Lock()
	defer h.Unlock()

	if h.pool == nil {
		h.status = false
		return errors.New("poll is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer func() {
		cancel()
		ctx.Done()
	}()

	conn, err := h.pool.Acquire(ctx)
	defer func() {
		if conn != nil {
			conn.Release()
		}
	}()

	if conn == nil || err != nil {
		h.status = false
		return err
	}
	h.status = true

	return nil
}

func (h *dbPoolHealth) checkStatus() {
	for {
		time.Sleep(2 * time.Second)
		err := dbPoolHealthInstance.dbPing()
		if err != nil {
			h.log.Warn(
				"db health checkStatus ",
				zap.String("error", fmt.Sprintf("%v", err)),
			)
		}
	}
}
