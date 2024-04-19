/*
 * This file was last modified at 2024-04-19 17:12 by Victor N. Skurikhin.
 * order_service.go
 * $Id$
 */

package orders

import (
	"context"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/vskurikhin/gophermart/internal/domain/dao"
	"github.com/vskurikhin/gophermart/internal/domain/entity"
	"github.com/vskurikhin/gophermart/internal/handlers"
	"github.com/vskurikhin/gophermart/internal/logger"
	"github.com/vskurikhin/gophermart/internal/model"
	"github.com/vskurikhin/gophermart/internal/storage"
	"github.com/vskurikhin/gophermart/internal/utils"
	"go.uber.org/zap"
	"net/http"
)

type OrderService interface {
	Number(login, number string) handlers.Result
	Orders(login string) handlers.Result
}

type service struct {
	ctx   context.Context
	log   *zap.Logger
	store *storage.PgsStorage
}

func newService(ctx context.Context) *service {
	return &service{
		ctx:   ctx,
		log:   logger.Get(),
		store: storage.NewPgsStorage(),
	}
}

func (s *service) Number(login, number string) handlers.Result {

	const funcName = "service.Number"
	defer utils.TraceInOut(s.ctx, funcName, "%s, %s", login, number)()

	if !checkLuhn(number) {
		return handlers.ResultErrorBadFormatNumber()
	}

	orders := dao.Orders(s.store.WithContext(s.ctx))
	order := entity.NewOrder(login, number)
	_, err := orders.Insert(order)

	if pe, ok := err.(*pgconn.PgError); ok {
		switch {
		case isIntegrityConstraintViolationOrderPkey(pe):
			return handlers.ResultErrorOrderByUserAlreadyLoaded()
		case isIntegrityConstraintViolationOrderNumberKey(pe):
			return handlers.ResultErrorOrderOtherAlreadyLoaded()
		}
		return handlers.ResultInternalError()
	}
	return handlers.NewResultString(number, http.StatusCreated)
}

func checkLuhn(number string) bool {

	var sum int
	parity := len(number) % 2

	for i := 0; i < len(number); i++ {
		digit := int(number[i] - '0')

		if i%2 == parity {
			digit *= 2

			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}
	return sum%10 == 0
}

func (s *service) Orders(login string) handlers.Result {

	const funcName = "service.Orders"
	defer utils.TraceInOut(s.ctx, funcName, "%s", login)()

	do := dao.Orders(s.store.WithContext(s.ctx))
	orders, err := do.GetAllOrdersForLogin(login)
	if err != nil {
		return handlers.ResultInternalError()
	}

	result := make(model.Orders, 0)
	for _, order := range orders {
		result = append(result, model.Order{
			Number:     order.Number(),
			UploadedAt: model.Time{Time: *order.UploadedAt()},
		})
	}

	return handlers.NewResultAny(result, http.StatusOK)
}

func isIntegrityConstraintViolationOrderPkey(pe *pgconn.PgError) bool {
	return pgerrcode.IsIntegrityConstraintViolation(pe.Code) &&
		pe.ConstraintName == "order_pkey"
}

func isIntegrityConstraintViolationOrderNumberKey(pe *pgconn.PgError) bool {
	return pgerrcode.IsIntegrityConstraintViolation(pe.Code) &&
		pe.ConstraintName == "order_number_key"
}