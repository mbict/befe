package jwtoken

import (
	"github.com/coreos/go-oidc/v3/oidc"
	"time"
)

func New(token *oidc.IDToken) *Token {
	claims := map[string]interface{}{}
	_ = token.Claims(&claims)

	return &Token{
		IDToken: token,
		claims:  claims,
	}
}

type Token struct {
	*oidc.IDToken
	claims map[string]interface{}
}

func (t *Token) Get(name string) (interface{}, bool) {
	v, ok := t.claims[name]
	return v, ok
}

func (t *Token) IsExpired() bool {
	return t.Expiry.Before(time.Now())
}
