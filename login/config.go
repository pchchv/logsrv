package login

import (
	"crypto/rand"
	"encoding/hex"
)

var jwtDefaultSecret string

func init() {
	var err error
	jwtDefaultSecret, err = randStringBytes(64)
	if err != nil {
		panic(err)
	}
}

func randStringBytes(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
