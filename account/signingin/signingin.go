package signingin

import (
	"log"

	"github.com/egoholic/editor/lib/pwd"
)

type (
	Value struct {
		logger      *log.Logger
		accessToken string
	}
	AccessTokenProvider interface {
		AccessToken(login, passwords string) (string, error)
	}
)

func New(l *log.Logger, atp AccessTokenProvider, login, password []byte) (*Value, error) {
	ep, err := pwd.Encrypt(password, login)
	if err != nil {
		return nil, err
	}
	token, err := atp.AccessToken(string(login), string(ep))
	if err != nil {
		return nil, err
	}
	return &Value{accessToken: token}, nil
}

func (v *Value) AccessToken() string {
	return v.accessToken
}
