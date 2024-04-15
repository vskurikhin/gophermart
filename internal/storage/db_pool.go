/*
 * This file was last modified at 2024-04-15 12:04 by Victor N. Skurikhin.
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

func DBPool() DB {
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
	_, err = pool.Acquire(context.TODO())

	if err != nil {
		panic(err)
	}
	log.Debug("Acquire pool Ok")

	return pool
}

func (this *dbPoolHealth) DBPool() (*pgxpool.Pool, bool) {
	this.RLock()
	defer this.RUnlock()
	return this.pool, this.status
}

func (this *dbPoolHealth) dbPing() error {
	this.Lock()
	defer this.Unlock()

	if this.pool == nil {
		this.status = false
		return errors.New("poll is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer func() {
		cancel()
		ctx.Done()
	}()

	conn, err := this.pool.Acquire(ctx)
	defer func() {
		if conn != nil {
			conn.Release()
		}
	}()

	if conn == nil || err != nil {
		this.status = false
		return err
	}
	this.status = true

	return nil
}

func (this *dbPoolHealth) checkStatus() {
	for {
		time.Sleep(2 * time.Second)
		err := dbPoolHealthInstance.dbPing()
		if err != nil {
			this.log.Warn(
				"db health checkStatus ",
				zap.String("error", fmt.Sprintf("%v", err)),
			)
		}
	}
}
