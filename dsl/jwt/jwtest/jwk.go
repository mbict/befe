package jwtest

import (
	"fmt"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
)

const keyId = "123450"

func PublicJwk() jwk.Key {
	key, err := jwk.New(privateKey.PublicKey)
	if err != nil {
		panic(fmt.Sprintf("failed to create public jwk key: %s\n", err))
	}
	key.Set(jwk.AlgorithmKey, jwa.RS256)
	key.Set(jwk.KeyIDKey, keyId)

	return key
}

func PrivateJwk() jwk.Key {
	key, err := jwk.New(privateKey)
	if err != nil {
		panic(fmt.Sprintf("failed to create private jwk key: %s\n", err))
	}
	key.Set(jwk.AlgorithmKey, jwa.RS256)
	key.Set(jwk.KeyIDKey, keyId)

	return key
}

func PublicJwkKeys() jwk.Set {
	set := jwk.NewSet()
	set.Add(PublicJwk())
	return set
}

func PrivateJwkKeys() jwk.Set {
	set := jwk.NewSet()
	set.Add(PrivateJwk())
	set.Add(PublicJwk())
	return set
}
