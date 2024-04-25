/*
 * This file was last modified at 2024-05-07 14:45 by Victor N. Skurikhin.
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

type Job struct {
	number    int
	requestID string
}

type Workers interface {
	Jobs() chan<- Job
}

type workers struct {
	jobs chan<- Job
}

func (w *workers) Jobs() chan<- Job {
	return w.jobs
}

var once = new(sync.Once)
var instance *workers

func GetWorkers() Workers {

	once.Do(func() {
		cfg := env.GetConfig()
		log := logger.Get()
		jobs := make(chan Job, cfg.RateLimit())

		w := new(workers)
		w.jobs = jobs

		for n := 1; n <= runtime.NumCPU(); n++ {
			go worker(log, n, jobs)
		}
		instance = w
	})
	return instance
}

func NewJob(number int, requestID string) Job {
	return Job{number: number, requestID: requestID}
}

func (j *Job) Number() int {
	return j.number
}

func (j *Job) RequestID() string {
	return j.requestID
}

func worker(log *zap.Logger, w int, jobs <-chan Job) {

	for job := range jobs {

		srv := newService(storage.NewPgsStorage(), job.RequestID())

		log.Debug("jobs",
			zap.String("number", fmt.Sprintf("%d", job.Number())),
			zap.String("worker", fmt.Sprintf("%d", w)),
		)
		err := srv.GetNumber(job.Number())

		if err != nil {
			log.Debug(
				"jobs",
				utils.LogCtxReasonErrFields(
					srv.Context(),
					fmt.Sprintf("number: %d, worker: %d", job.Number(), w), err,
				)...,
			)
		}
	}
}
