package token

import (
	"github.com/lestrrat-go/jwx/jwt"
)

type JwtDecoder interface {
	Decode(token []byte) (*JwtToken, error)
}

func NewJwtDecoder() JwtDecoder {
	return &jwtDecoder{}
}

type jwtDecoder struct {
}

func (d *jwtDecoder) Decode(jwtToken []byte) (*JwtToken, error) {
	token, err := jwt.Parse(jwtToken)
	if err != nil {
		return nil, err
	}
	return &JwtToken{
		Token: token,
	}, nil
}
