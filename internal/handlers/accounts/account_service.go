/*
 * This file was last modified at 2024-04-19 20:15 by Victor N. Skurikhin.
 * account_service.go
 * $Id$
 */

package accounts

import (
	"context"
	"github.com/vskurikhin/gophermart/internal/domain/dao"
	"github.com/vskurikhin/gophermart/internal/handlers"
	"github.com/vskurikhin/gophermart/internal/logger"
	"github.com/vskurikhin/gophermart/internal/model"
	"github.com/vskurikhin/gophermart/internal/storage"
	"github.com/vskurikhin/gophermart/internal/utils"
	"go.uber.org/zap"
	"net/http"
)

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

	const funcName = "service.Balance"
	defer utils.TraceInOut(s.ctx, funcName, "%s", login)()

	db := dao.Balances(s.store.WithContext(s.ctx))
	b, sum, err := db.GetBalanceWithdraw(login)

	if err != nil && !utils.IsErrNoRowsInResultSet(err) {
		return handlers.ResultInternalError()
	} else if utils.IsErrNoRowsInResultSet(err) {
		s.log.Debug(balanceMsg, utils.LogCtxReasonErrFields(s.ctx, err.Error(), handlers.ErrBalanceNotSet)...)
	}
	balance := model.NewBalanceBigFloat(b.Balance(), *sum)

	return handlers.NewResultAny(balance, http.StatusOK)
}
