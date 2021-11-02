package token

type JwtVerifier interface {
	Verify(jwtToken []byte) (valid bool, token *JwtToken, err error)
}

type jwtHS256 struct {
}

func (v *jwtHS256) Verify(jwtToken []byte) (valid bool, token *JwtToken, err error) {
	panic("implement me")
}

func NewJwtVerifierHS256(secret []byte) JwtVerifier {
	return &jwtHS256{}
}
