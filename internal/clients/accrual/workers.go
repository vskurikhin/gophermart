/*
 * This file was last modified at 2024-04-20 17:09 by Victor N. Skurikhin.
 * workers.go
 * $Id$
 */

package accrual

import (
	"fmt"
	"github.com/vskurikhin/gophermart/internal/env"
	"github.com/vskurikhin/gophermart/internal/logger"
	"github.com/vskurikhin/gophermart/internal/storage"
	"github.com/vskurikhin/gophermart/internal/utils"
	"go.uber.org/zap"
	"runtime"
	"sync"
)

type Workers interface {
	Jobs() chan<- int
}

type workers struct {
	jobs chan<- int
}

func (w *workers) Jobs() chan<- int {
	return w.jobs
}

var once = new(sync.Once)
var instance *workers

func GetWorkers() Workers {

	once.Do(func() {
		cfg := env.GetConfig()
		jobs := make(chan int, cfg.RateLimit())
		log := logger.Get()

		w := new(workers)
		w.jobs = jobs

		for n := 1; n <= runtime.NumCPU(); n++ {
			go worker(log, n, jobs)
		}
		instance = w
	})
	return instance
}

func worker(log *zap.Logger, w int, jobs <-chan int) {

	srv := newService(storage.NewPgsStorage())

	for i := range jobs {

		log.Debug("jobs",
			zap.String("number", fmt.Sprintf("%d", i)),
			zap.String("worker", fmt.Sprintf("%d", w)),
		)
		err := srv.GetNumber(i)

		if err != nil {
			log.Debug(
				"jobs",
				utils.LogCtxReasonErrFields(srv.Context(), fmt.Sprintf("number: %d, worker: %d", i, w), err)...,
			)
		}
	}
}
