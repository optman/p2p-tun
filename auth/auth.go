package auth

import (
	"bytes"
	"crypto/sha256"
	"io"
)

type Authenticator struct {
	key []byte
}

func NewAuthenticator(secret string) *Authenticator {
	h := sha256.New()
	h.Write([]byte(secret))

	return &Authenticator{
		key: h.Sum(nil),
	}
}

func (au *Authenticator) Write(w io.Writer) error {
	_, err := w.Write(au.key)
	return err
}

func (au *Authenticator) Read(r io.Reader) (bool, error) {
	buf := make([]byte, 32)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return false, err
	}

	return bytes.Compare(au.key, buf) == 0, nil
}
