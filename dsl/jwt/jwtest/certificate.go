package jwtest

import (
	"crypto/rand"
	"crypto/rsa"
)

var privateKey *rsa.PrivateKey

func init() {
	var err error
	privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic("cannot generate private sign key for testing: " + err.Error())
	}
}
