package token

import (
	"context"
	"github.com/lestrrat-go/backoff/v2"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"time"
)

type jwkVerifier struct {
	issuer        string
	keysetFetcher *jwk.AutoRefresh
}

func (v *jwkVerifier) Verify(jwtToken []byte) (bool, *JwtToken, error) {

	if jwtToken == nil || len(jwtToken) == 0 {
		return false, nil, nil
	}

	keyset, err := v.fetchJwkKeys()
	if err != nil {
		return false, nil, err
	}

	token, err := jwt.Parse(jwtToken,
		jwt.WithKeySet(keyset),
		jwt.WithValidate(true))

	if err != nil {
		//log.Default().Println(err)
		return false, nil, err
	}

	return true, &JwtToken{
		Token: token,
	}, nil
}

func (v *jwkVerifier) fetchJwkKeys() (jwk.Set, error) {
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	return v.keysetFetcher.Fetch(ctx, v.issuer)
}

func NewJwkVerifier(ctx context.Context, issuer string) JwtVerifier {
	jwkFetcher := jwk.NewAutoRefresh(ctx)
	jwkFetcher.Configure(issuer,
		jwk.WithMinRefreshInterval(15*time.Minute),
		jwk.WithFetchBackoff(backoff.Constant(backoff.WithInterval(time.Minute))),
	)

	return &jwkVerifier{
		keysetFetcher: jwkFetcher,
		issuer:        issuer,
	}
}
