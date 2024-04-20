/*
 * This file was last modified at 2024-04-20 17:09 by Victor N. Skurikhin.
 * order_service.go
 * $Id$
 */

package orders

import (
	"context"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/vskurikhin/gophermart/internal/clients/accrual"
	"github.com/vskurikhin/gophermart/internal/domain/dao"
	"github.com/vskurikhin/gophermart/internal/domain/entity"
	"github.com/vskurikhin/gophermart/internal/handlers"
	"github.com/vskurikhin/gophermart/internal/logger"
	"github.com/vskurikhin/gophermart/internal/model"
	"github.com/vskurikhin/gophermart/internal/storage"
	"github.com/vskurikhin/gophermart/internal/utils"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

const orderServiceMsg = "order service"

type OrderService interface {
	Number(login, number string) handlers.Result
	Orders(login string) handlers.Result
}

type service struct {
	ctx     context.Context
	log     *zap.Logger
	store   storage.Storage
	workers accrual.Workers
}

func newService(ctx context.Context, store storage.Storage, workers accrual.Workers) *service {
	return &service{
		ctx:     ctx,
		log:     logger.Get(),
		store:   store,
		workers: workers,
	}
}

func (s *service) Number(login, number string) handlers.Result {

	const funcName = "service.Number"
	defer utils.TraceInOut(s.ctx, funcName, "%s, %s", login, number)()

	if !utils.CheckLuhn(number) {
		return handlers.ResultErrorBadFormatNumber()
	}
	job, err := strconv.Atoi(number)

	if err != nil {
		return handlers.ResultErrorBadFormatNumber()
	}

	orders := dao.Orders(s.store.WithContext(s.ctx))
	order := entity.NewOrder(login, number)
	_, err = orders.Insert(order)

	if pe, ok := err.(*pgconn.PgError); ok {
		switch {
		case isIntegrityConstraintViolationOrdersPkey(pe):
			return handlers.ResultErrorOrderByUserAlreadyLoaded()
		case isIntegrityConstraintViolationOrdersNumberKey(pe):
			return handlers.ResultErrorOrderOtherAlreadyLoaded()
		}
		s.log.Debug(orderServiceMsg, utils.LogCtxReasonErrFields(s.ctx, "insert", err)...)
		return handlers.ResultInternalError()
	}
	jobs := s.workers.Jobs()
	jobs <- job

	return handlers.NewResultString(number, http.StatusAccepted)
}

func (s *service) Orders(login string) handlers.Result {

	const funcName = "service.Orders"
	defer utils.TraceInOut(s.ctx, funcName, "%s", login)()

	do := dao.Orders(s.store.WithContext(s.ctx))
	orders, err := do.GetAllOrdersForLogin(login)

	if err != nil {
		s.log.Debug(orderServiceMsg, utils.LogCtxReasonErrFields(s.ctx, "get order", err)...)
		return handlers.ResultInternalError()
	}
	result := make(model.Orders, 0)

	for _, order := range orders {

		var accuracy *model.Float

		if order.Accrual() != nil {
			a := order.Accrual()
			f, _ := a.Float64()
			modelFloat := model.Float(f)
			accuracy = &modelFloat
		}
		result = append(result, model.Order{
			Number:     order.Number(),
			Status:     order.Status(),
			Accrual:    accuracy,
			UploadedAt: model.Time{Time: *order.UploadedAt()},
		})
	}

	return handlers.NewResultAny(result, http.StatusOK)
}

func isIntegrityConstraintViolationOrdersPkey(pe *pgconn.PgError) bool {
	return pgerrcode.IsIntegrityConstraintViolation(pe.Code) &&
		pe.ConstraintName == "orders_pkey"
}

func isIntegrityConstraintViolationOrdersNumberKey(pe *pgconn.PgError) bool {
	return pgerrcode.IsIntegrityConstraintViolation(pe.Code) &&
		pe.ConstraintName == "orders_number_key"
}
