package signingin

import (
	"errors"
	"log"
	"time"

	"github.com/egoholic/editor/lib/pwd"
)

type (
	Signin struct {
		Login     string
		CreatedAt time.Time
	}
	Value struct {
		logger      *log.Logger
		accessToken string
	}
	AccountAuthenticator interface {
		IsAuthenticated(login, password string) (ok bool, err error)
	}
	SigninSaver interface {
		Save(*Signin) (string, error)
	}
)

func New(l *log.Logger, aa AccountAuthenticator, sis SigninSaver, login, password string) (*Value, error) {
	var (
		err   error
		token string
	)
	ep, err := pwd.Encrypt([]byte(password), []byte(login))
	if err != nil {
		return nil, err
	}
	ok, err := aa.IsAuthenticated(login, string(ep))
	if err != nil {
		return nil, err
	}
	if !ok {
		err = errors.New("wrong login or password")
		return nil, err
	}
	signin := &Signin{
		Login:     string(login),
		CreatedAt: time.Now().UTC(),
	}
	token, err = sis.Save(signin)
	if err != nil {
		return nil, err
	}
	return &Value{
		logger:      l,
		accessToken: token,
	}, nil
}

func (v *Value) AccessToken() string {
	return v.accessToken
}
