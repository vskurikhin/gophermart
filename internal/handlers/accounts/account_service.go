/*
 * This file was last modified at 2024-04-20 17:09 by Victor N. Skurikhin.
 * account_service.go
 * $Id$
 */

package accounts

import (
	"context"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/vskurikhin/gophermart/internal/domain/dao"
	"github.com/vskurikhin/gophermart/internal/domain/entity"
	"github.com/vskurikhin/gophermart/internal/domain/transaction"
	"github.com/vskurikhin/gophermart/internal/handlers"
	"github.com/vskurikhin/gophermart/internal/logger"
	"github.com/vskurikhin/gophermart/internal/model"
	"github.com/vskurikhin/gophermart/internal/storage"
	"github.com/vskurikhin/gophermart/internal/utils"
	"go.uber.org/zap"
	"math/big"
	"net/http"
)

const accountServiceMsg = "account service"

type AccountService interface {
	Balance(login string) handlers.Result
}

type service struct {
	ctx   context.Context
	log   *zap.Logger
	store storage.Storage
}

func newService(ctx context.Context, store storage.Storage) *service {
	return &service{
		ctx:   ctx,
		log:   logger.Get(),
		store: store,
	}
}

func (s *service) Balance(login string) handlers.Result {

	const funcName = "service.Current"
	defer utils.TraceInOut(s.ctx, funcName, "%s", login)()

	db := dao.Balances(s.store.WithContext(s.ctx))
	b, sum, err := db.GetBalanceWithdraw(login)

	if err != nil && !utils.IsErrNoRowsInResultSet(err) {
		s.log.Debug(accountServiceMsg, utils.LogCtxReasonErrFields(s.ctx, "get balance", err)...)
		return handlers.ResultInternalError()
	} else if utils.IsErrNoRowsInResultSet(err) {
		s.log.Debug(accountServiceMsg, utils.LogCtxReasonErrFields(s.ctx, err.Error(), handlers.ErrBalanceNotSet)...)
	}
	if sum == nil {
		sum = big.NewFloat(0)
	}
	balance := model.NewBalanceBigFloat(b.Current(), *sum)

	return handlers.NewResultAny(balance, http.StatusOK)
}

func (s *service) Withdraw(login string, modelWithdraw *model.Withdraw) handlers.Result {

	const funcName = "service.Withdraw"
	defer utils.TraceInOut(s.ctx, funcName, "%s, %v", login, modelWithdraw)()

	if !utils.CheckLuhn(modelWithdraw.Order) {
		return handlers.ResultErrorBadFormatNumber()
	}
	db := dao.Balances(s.store.WithContext(s.ctx))
	entityBalance, err := db.GetBalance(login)

	if err != nil {
		return handlers.ResultInternalError()
	}
	balance := entityBalance.Current()
	sum := big.NewFloat(modelWithdraw.GetSum())

	if sum.Cmp(&balance) > 0 {
		return handlers.ResultErrorPaymentRequired()
	}
	entityWithdraw := entity.NewWithdraw(login, modelWithdraw.Order, *sum)
	tbw := transaction.BalanceWithdraw(s.store.WithContext(s.ctx))
	err = tbw.TransactionWithdraw(entityBalance, entityWithdraw)

	if pe, ok := err.(*pgconn.PgError); ok {
		switch {
		case isIntegrityConstraintViolationWithdrawPkey(pe):
			return handlers.ResultErrorOrderByUserAlreadyLoaded()
		}
	}
	if err != nil {
		s.log.Debug(accountServiceMsg, utils.LogCtxReasonErrFields(s.ctx, "transaction withdraw", err)...)
		return handlers.ResultInternalError()
	}
	b := entityWithdraw.Sum()
	withdraw := model.NewWithdraw(entityWithdraw.Number(), &b)

	return handlers.NewResultAny(withdraw, http.StatusOK)
}

func (s *service) Withdrawals(login string) handlers.Result {

	const funcName = "service.Withdrawals"
	defer utils.TraceInOut(s.ctx, funcName, "%s", login)()

	dw := dao.Withdrawals(s.store.WithContext(s.ctx))
	withdrawals, err := dw.GetAllWithdrawalsByLogin(login)
	if err != nil {
		return handlers.ResultInternalError()
	}

	result := make(model.Withdrawals, 0)
	for _, withdraw := range withdrawals {
		bs := withdraw.Sum()
		sum, _ := bs.Float64()
		var processedAt *model.Time
		if withdraw.ProcessedAt() == nil {
			processedAt = &model.Time{Time: withdraw.CreatedAt()}
		} else {
			processedAt = &model.Time{Time: *withdraw.ProcessedAt()}
		}
		result = append(result, model.Withdraw{
			Order:       withdraw.Number(),
			Sum:         model.Float(sum),
			ProcessedAt: processedAt,
		})
	}
	if len(result) == 0 {
		return handlers.NewResultString("[]", http.StatusNoContent)
	}
	return handlers.NewResultAny(result, http.StatusOK)
}

func isIntegrityConstraintViolationWithdrawPkey(pe *pgconn.PgError) bool {
	return pgerrcode.IsIntegrityConstraintViolation(pe.Code) &&
		pe.ConstraintName == "withdraw_pkey"
}
