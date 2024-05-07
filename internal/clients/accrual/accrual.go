/*
 * This file was last modified at 2024-05-07 14:41 by Victor N. Skurikhin.
 * accrual.go
 * $Id$
 */

package accrual

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vskurikhin/gophermart/internal/domain/dao"
	"github.com/vskurikhin/gophermart/internal/domain/transaction"
	"github.com/vskurikhin/gophermart/internal/env"
	"github.com/vskurikhin/gophermart/internal/handlers"
	"github.com/vskurikhin/gophermart/internal/logger"
	"github.com/vskurikhin/gophermart/internal/model"
	"github.com/vskurikhin/gophermart/internal/storage"
	"github.com/vskurikhin/gophermart/internal/utils"
	"go.uber.org/zap"
	"math/big"
	"net/http"
	"time"
)

const (
	increase  = 2
	tries     = 3
	Invalid   = "INVALID"
	Processed = "PROCESSED"
)

type AccrualsService interface {
	Context() context.Context
	GetNumber(number int) error
}

type service struct {
	address string
	client  *http.Client
	ctx     context.Context
	log     *zap.Logger
	store   storage.Storage
}

func newService(store storage.Storage, id string) *service {

	cfg := env.GetConfig()
	ctx := context.WithValue(context.Background(), middleware.RequestIDKey, id)

	return &service{
		address: cfg.AccrualSystemAddress(),
		client:  &http.Client{},
		ctx:     ctx,
		log:     logger.Get(),
		store:   store.WithContext(ctx),
	}
}

func (s *service) Context() context.Context {
	return s.ctx
}

func (s *service) GetNumber(number int) error {

	utils.TraceInOut(s.ctx, "AccrualsService.GetNumber", "%d", number)

	var err error
	var done bool

	for !done {
		done, err = s.getNumber(number)
		for i := 1; err != nil && i < tries*increase; i += increase {
			time.Sleep(time.Duration(i) * time.Second)
			s.log.Warn("jobs", utils.LogCtxReasonErrFields(s.ctx, "retry get accrual", err)...)
			done, err = s.getNumber(number)
		}
		if err != nil {
			return err
		}
	}
	s.log.Info(
		"GetNumber",
		zap.String("reqId", middleware.GetReqID(s.ctx)),
		zap.String("status", "done"),
		zap.Int("number", number),
	)
	return nil
}

//goland:noinspection GoUnhandledErrorResult
func (s *service) getNumber(number int) (bool, error) {

	url := fmt.Sprintf("http://%s/api/orders/%d", s.address, number)
	response, err := http.Get(url)
	defer func() {
		if response != nil {
			response.Body.Close()
		}
	}()

	if err != nil {
		s.log.Debug(
			"jobs",
			utils.LogCtxReasonErrFields(s.ctx, fmt.Sprintf("number: %d", number), err)...,
		)
		return false, err
	}
	if response.StatusCode == http.StatusOK {
		accrual, err := model.UnmarshalAccrualFromReader(response.Body)
		if err != nil {
			s.log.Debug(
				"jobs",
				utils.LogCtxReasonErrFields(s.ctx, fmt.Sprintf("response: %v", response), err)...,
			)
			return false, err
		}
		do := dao.Orders(s.store.WithContext(s.ctx))
		order, err := do.GetOrderByNumber(accrual.Order)
		if err != nil {
			s.log.Debug(
				"jobs",
				utils.LogCtxReasonErrFields(s.ctx, fmt.Sprintf("response: %v", response), err)...,
			)
			return false, err
		}
		db := dao.Balances(s.store.WithContext(s.ctx))
		balance, err := db.GetBalance(order.Login())
		if err != nil {
			s.log.Debug(
				"jobs",
				utils.LogCtxReasonErrFields(s.ctx, fmt.Sprintf("response: %v", response), err)...,
			)
			return false, err
		}
		bo := transaction.BalanceOrder(s.store.WithContext(s.ctx))
		sum := big.NewFloat(accrual.GetAccrual())
		order.SetAccrual(sum)
		err = bo.TransactionAccrual(balance, order, accrual.Status)
		if err != nil {
			s.log.Debug(
				"jobs",
				utils.LogCtxReasonErrFields(s.ctx, fmt.Sprintf("response: %v", response), err)...,
			)
			return false, err
		}
		return accrual.Status == Invalid || accrual.Status == Processed, nil
	}
	return false, handlers.ErrBalanceNotSet
}
