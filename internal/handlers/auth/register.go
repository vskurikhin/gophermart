/*
 * This file was last modified at 2024-04-15 16:54 by Victor N. Skurikhin.
 * register.go
 * $Id$
 */

package auth

import (
	"github.com/go-chi/render"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/vskurikhin/gophermart/internal/domain/dao"
	"github.com/vskurikhin/gophermart/internal/domain/entity"
	"github.com/vskurikhin/gophermart/internal/handlers"
	"github.com/vskurikhin/gophermart/internal/model"
	"github.com/vskurikhin/gophermart/internal/storage"
	"github.com/vskurikhin/gophermart/internal/utils"
	"go.uber.org/zap"
	"net/http"
)

const regMsg = "register"

type register struct {
	log   *zap.Logger
	store *storage.PgsStorage
}

type registerError struct {
	err    error
	status int
}

func newRegister(log *zap.Logger, store *storage.PgsStorage) *register {
	return &register{log: log, store: store}
}

func (r *register) Handle(response http.ResponseWriter, request *http.Request) {

	userRegister, err := model.UnmarshalFromReader(request.Body)
	_, registerError := r.doReg(userRegister, err)

	if registerError != nil {
		render.Status(request, registerError.status)
		//goland:noinspection GoUnhandledErrorResult
		render.Render(response, request, model.Error(registerError.err))
		return
	}
	if err := userRegister.MarshalToWriter(response); err != nil {
		panic(err)
	}
}

func (r *register) doReg(userRegister *model.UserRegister, err error) (login *string, re *registerError) {

	login, err = r.doRegister(userRegister, err)

	if pe, ok := err.(*pgconn.PgError); ok && isIntegrityConstraintViolation(pe) {
		re = &registerError{
			err:    handlers.ErrStatusConflict,
			status: http.StatusConflict,
		}
	} else if err != nil {
		re = &registerError{
			err:    handlers.ErrBadRequest,
			status: http.StatusBadRequest,
		}
	}
	return
}

func (r *register) doRegister(userRegister *model.UserRegister, err error) (*string, error) {

	ctx := utils.NewUUIDContext()

	if err != nil {
		r.log.Debug(regMsg, utils.LogCtxRecoverFields(ctx, err)...)
		return nil, err
	}
	hashed, err := utils.HashPassword(userRegister.Password)

	if err != nil {
		r.log.Debug(regMsg, utils.LogCtxRecoverFields(ctx, err)...)
		return nil, err
	}
	user := entity.NewUser(userRegister.Login, &hashed)
	users := dao.Users(r.store.WithContext(ctx))
	user, err = users.Insert(user)
	if err != nil {
		r.log.Debug(regMsg, utils.LogCtxRecoverFields(ctx, err)...)
		return nil, err
	}
	login := user.Login()

	return &login, nil
}

func isIntegrityConstraintViolation(pe *pgconn.PgError) bool {
	return pgerrcode.IsIntegrityConstraintViolation(pe.Code) &&
		pe.ConstraintName == "user_pkey"
}
