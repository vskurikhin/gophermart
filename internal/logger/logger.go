/*
 * This file was last modified at 2024-04-13 19:03 by Victor N. Skurikhin.
 * logger.go
 * $Id$
 */

package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
)

var log *zap.Logger
var Development bool
var once = new(sync.Once)

func Get() *zap.Logger {
	once.Do(func() {
		if Development {
			config := zap.NewDevelopmentConfig()
			log = zap.Must(config.Build())
		} else {
			config := zap.NewProductionConfig()
			config.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
			log = zap.Must(config.Build())
		}
	})
	//goland:noinspection GoUnhandledErrorResult
	defer log.Sync()
	return log
}
