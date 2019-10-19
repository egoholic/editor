package pwd

import (
	"golang.org/x/crypto/scrypt"
)

type EncryptedPassword []byte

const (
	cost   = 65536
	r      = 16
	p      = 4
	keyLen = 64
)

func Encrypt(password, salt []byte) ([]byte, error) {
	return scrypt.Key(password, salt, cost, r, p, keyLen)
}

func (ep EncryptedPassword) Compare(password, salt []byte) bool {
	ep2, err := Encrypt(password, salt)
	if err != nil {
		return false
	}
	epBytes := []byte(ep)
	ep2Bytes := []byte(ep2)
	for i, b := range epBytes {
		if b != ep2Bytes[i] {
			return false
		}
	}
	return true
}
