package auth

import (
	"context"
	jwtoken "culc/iternal/lib/jwt"
	"culc/iternal/model"
	"culc/iternal/postg"
	"errors"
	"fmt"
	"time"
)

type Auth struct {
	usrSaver Storage
	tokentl  time.Duration
}
type Storage interface {
	SaveUser(ctx context.Context, email string, password string) (uid uint, err error)
	UserProvider(ctx context.Context, email string) (model.User, error)
	SaveEx(ctx context.Context, ex string, uid float64) (int64, []string, error)
	UpdateSubEx(ctx context.Context, idex int64, userid float64, subEx []string) error
	UpdateExTimeReady(ctx context.Context, idex int64, userid float64) error
	GetExHistory(ctx context.Context, userid float64) ([]model.Expression, error)
	GetUnREadyEx() ([]model.Expression, error)
	ErrorinEx(float64, int64) error
}

func New(usrSaver Storage, tokentl time.Duration) *Auth {
	return &Auth{
		usrSaver: usrSaver,
		tokentl:  tokentl,
	}
}
func (a *Auth) Login(ctx context.Context, email string, password string) (string, error) {
	user, err := a.usrSaver.UserProvider(ctx, email)
	if err != nil {
		if errors.Is(err, postg.ErrUserExist) {
			return "", fmt.Errorf("%e", err)
		}
	}

	if user.Password != password {

		return "", fmt.Errorf("error in login ")
	}
	token, err := jwtoken.NewToken(user, a.tokentl)
	if err != nil {
		return "", fmt.Errorf("error in token")
	}
	return token, nil

}
func (a *Auth) Register(ctx context.Context, email string, password string) (uint, error) {
	id, err := a.usrSaver.SaveUser(ctx, email, password)
	if err != nil {
		return 0, fmt.Errorf("error in registr %e", err)
	}
	return id, nil
}
func (a *Auth) SaveExpression(ctx context.Context, ex string, uid float64) (int64, []string, error) {

	id, subEx, err := a.usrSaver.SaveEx(ctx, ex, uid)
	if err != nil {
		return 0, nil, fmt.Errorf("error add ex to db")
	}

	return id, subEx, nil
}
func (a *Auth) UpdateSubEx(ctx context.Context, idex int64, iduser float64, subEx []string) error {
	err := a.usrSaver.UpdateSubEx(ctx, idex, iduser, subEx)
	if err != nil {
		return err
	}
	return nil
}
func (a *Auth) GetExHistory(ctx context.Context, userid float64) ([]model.Expression, error) {
	Ex, err := a.usrSaver.GetExHistory(ctx, userid)
	if err != nil {
		return nil, err
	}
	return Ex, err

}
func (a *Auth) UpdateExTimeReady(ctx context.Context, idex int64, userid float64) error {
	err := a.usrSaver.UpdateExTimeReady(ctx, idex, userid)
	if err != nil {
		return err
	}
	return nil
}
func (a *Auth) GetUnREadyEx() ([]model.Expression, error) {
	ex, err := a.usrSaver.GetUnREadyEx()
	if err != nil {
		return nil, err
	}
	return ex, nil
}
func (a *Auth) ErrorinEx(uid float64, idex int64) error {
	err := a.usrSaver.ErrorinEx(uid, idex)
	if err != nil {
		return err
	}
	return nil
}
