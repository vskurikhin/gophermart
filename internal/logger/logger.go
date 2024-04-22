/*
 * This file was last modified at 2024-04-21 01:19 by Victor N. Skurikhin.
 * logger.go
 * $Id$
 */

package logger

import (
	"github.com/vskurikhin/gophermart/internal/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
)

var log *zap.Logger
var once = new(sync.Once)

func Get() *zap.Logger {
	once.Do(func() {
		cfg := env.GetConfig()
		if cfg.DevelopmentLogger() {
			config := zap.NewDevelopmentConfig()
			log = zap.Must(config.Build())
		} else {
			config := zap.NewProductionConfig()
			config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
			log = zap.Must(config.Build())
		}
	})
	//goland:noinspection GoUnhandledErrorResult
	defer log.Sync()
	return log
}
