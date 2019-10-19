package signingin

import (
	"github.com/egoholic/editor/lib/pwd"
)

type (
	Value struct {
		accessToken string
	}
	AccessTokenProvider interface {
		AccessToken(login, passwords string) (string, error)
	}
)

func New(atp AccessTokenProvider, login, password []byte) (*Value, error) {
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
