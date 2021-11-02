package token

import "github.com/lestrrat-go/jwx/jwt"

type JwtToken struct {
	jwt.Token
}

func (t *JwtToken) IsExpired() bool {
	return false
}
