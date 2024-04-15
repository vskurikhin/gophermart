/*
 * This file was last modified at 2024-04-16 17:02 by Victor N. Skurikhin.
 * user_service.go
 * $Id$
 */

package auth

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

type UserService interface {
	Login(userRegister *model.User) handlers.Result
	Register(userRegister *model.User) handlers.Result
}

type userService struct {
	ctx   context.Context
	log   *zap.Logger
	store *storage.PgsStorage
}

func NewUserService(ctx context.Context) UserService {
	return &userService{
		ctx:   ctx,
		log:   logger.Get(),
		store: storage.NewPgsStorage(),
	}
}

func (u *userService) Login(user *model.User) handlers.Result {

	const funcName = "userService.Login"
	defer utils.TraceInOut(u.ctx, funcName, "%+v", *user)()

	token, err := u.login(user)

	if err != nil {
		return handlers.NewResultError(err, http.StatusUnauthorized)
	}
	return handlers.NewResultString(token, http.StatusOK)
}

func (u *userService) login(modelUser *model.User) (string, error) {

	users := dao.Users(u.store.WithContext(u.ctx))
	entityUser, err := users.GetUser(modelUser.Login)

	if err == nil && entityUser != nil && passwordOk(modelUser, entityUser) {

		login := entityUser.Login()
		token := utils.MakeToken(login)

		return token, nil
	}
	u.log.Debug(logMsg, utils.LogCtxRecoverFields(u.ctx, handlers.ErrBadUserPassword)...)

	return "", handlers.ErrBadUserPassword
}

func (u *userService) Register(modelUser *model.User) handlers.Result {

	const funcName = "userService.Register"
	defer utils.TraceInOut(u.ctx, funcName, "%+v", *modelUser)()

	token, err := u.register(modelUser)

	if pe, ok := err.(*pgconn.PgError); ok && isIntegrityConstraintViolation(pe) {
		return handlers.NewResultError(handlers.ErrStatusConflict, http.StatusConflict)
	} else if err != nil {
		return handlers.NewResultError(handlers.ErrBadRequest, http.StatusBadRequest)
	} else if token == "" {
		return handlers.NewResultError(handlers.ErrBadRequest, http.StatusInternalServerError)
	}
	return handlers.NewResultString(token, http.StatusOK)
}

func (u *userService) register(modelUser *model.User) (string, error) {

	hashed, err := utils.HashPassword(modelUser.Password)

	if err != nil {
		u.log.Debug(regMsg, utils.LogCtxRecoverFields(u.ctx, err)...)
		return "", err
	}
	entityUser := entity.NewUser(modelUser.Login, &hashed)
	users := dao.Users(u.store.WithContext(u.ctx))
	entityUser, err = users.Insert(entityUser)

	if err != nil {
		u.log.Debug(regMsg, utils.LogCtxRecoverFields(u.ctx, err)...)
		return "", err
	}
	login := entityUser.Login()
	token := utils.MakeToken(login)

	return token, nil
}

func isIntegrityConstraintViolation(pe *pgconn.PgError) bool {
	return pgerrcode.IsIntegrityConstraintViolation(pe.Code) &&
		pe.ConstraintName == "user_pkey"
}

func passwordOk(modelUser *model.User, entityUser *entity.User) bool {
	return entityUser.Password() != nil &&
		utils.CheckPasswordHash(modelUser.Password, *entityUser.Password())
}
