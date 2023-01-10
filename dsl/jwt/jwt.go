package jwt

import (
	"context"
	"github.com/coreos/go-oidc/v3/oidc"
	. "github.com/mbict/befe/dsl"
	"github.com/mbict/befe/dsl/jwt/jwtoken"
	. "github.com/mbict/befe/expr"
	"log"
	"net/http"
	"strings"
)

var stopHandler = func(rw http.ResponseWriter, r *http.Request) (bool, error) {
	return false, nil
}

type TokenFetcher func(r *http.Request) string

func TokenFromHeader() TokenFetcher {
	return func(r *http.Request) string {
		reqToken := r.Header.Get("Authorization")
		splitToken := strings.Split(reqToken, "Bearer ")
		if len(splitToken) != 2 {
			return ""
		}
		return splitToken[1]
	}
}

func TokenFromCookie(name string) TokenFetcher {
	return func(r *http.Request) string {
		c, err := r.Cookie(name)
		if err != nil {
			return ""
		}
		return c.Value
	}
}

type Jwk interface {
	Action

	WithAudience(audiences ...string) Jwk
	WithExpiredCheck() Jwk
	WithClaim(name string, values ...string) Jwk

	WhenExpired(actions ...Action) Jwk
	WhenDenied(actions ...Action) Jwk
	WhenNoToken(actions ...Action) Jwk

	TokenFrom(tokenFetchers ...TokenFetcher) Jwk
}

func JwkToken(issuer string) Jwk {

	return &jwtValidator{
		tokenFetchers:     []TokenFetcher{TokenFromHeader()},
		issuer:            issuer,
		withAudienceCheck: nil,
		withExpiredCheck:  false,
		withClaims:        map[string][]string{},
		hasClaims:         nil,
		expiredActions:    nil,
		deniedActions:     nil,
		noTokenActions:    nil,
	}
}

type jwtValidator struct {
	tokenFetchers []TokenFetcher
	issuer        string

	withAudienceCheck []string
	withExpiredCheck  bool
	withClaims        map[string][]string
	hasClaims         []string

	expiredActions Actions
	deniedActions  Actions
	noTokenActions Actions
}

func (j *jwtValidator) WithAudience(audiences ...string) Jwk {
	j.withAudienceCheck = append(j.withAudienceCheck, audiences...)
	return j
}

func (j *jwtValidator) WithExpiredCheck() Jwk {
	j.withExpiredCheck = true
	return j
}

func (j *jwtValidator) WithClaim(name string, values ...string) Jwk {
	j.withClaims[name] = append(j.withClaims[name], values...)
	return j
}

func (j *jwtValidator) WhenExpired(actions ...Action) Jwk {
	j.expiredActions = append(j.expiredActions, actions...)
	return j
}

func (j *jwtValidator) WhenDenied(actions ...Action) Jwk {
	j.deniedActions = append(j.deniedActions, actions...)
	return j
}

func (j *jwtValidator) WhenNoToken(actions ...Action) Jwk {
	j.noTokenActions = append(j.noTokenActions, actions...)
	return j
}

func (j *jwtValidator) TokenFrom(fetchers ...TokenFetcher) Jwk {
	j.tokenFetchers = fetchers
	return j
}

func (j *jwtValidator) BuildHandler(ctx context.Context, next Handler) Handler {
	jwtVerifier := oidc.NewVerifier(j.issuer, oidc.NewRemoteKeySet(ctx, j.issuer), &oidc.Config{
		ClientID:                   "",
		SupportedSigningAlgs:       []string{"RS256"},
		SkipClientIDCheck:          true,
		SkipExpiryCheck:            !j.withExpiredCheck,
		SkipIssuerCheck:            true, // <- make it configurable
		InsecureSkipSignatureCheck: false,
	})

	if j.deniedActions == nil {
		j.deniedActions = Actions{Deny()}
	}
	deniedHandler := j.deniedActions.BuildHandler(ctx, stopHandler)

	if j.expiredActions == nil {
		j.expiredActions = Actions{Deny()}
	}
	expiredHandler := j.expiredActions.BuildHandler(ctx, stopHandler)

	if j.noTokenActions == nil {
		j.noTokenActions = Actions{Unauthorized()}
	}
	noTokenHandler := j.noTokenActions.BuildHandler(ctx, stopHandler)

	conditions := Conditions{}
	if len(j.withAudienceCheck) > 0 {
		conditions = append(conditions, HasJwtAudience(j.withAudienceCheck...))
	}

	if len(j.withClaims) > 0 {
		for key, values := range j.withClaims {
			conditions = append(conditions, HasJwtClaim(key, values...))
		}
	}
	postCheck := conditions.BuildCondition(ctx)

	return func(rw http.ResponseWriter, r *http.Request) (bool, error) {
		var jwt string
		for _, fetcher := range j.tokenFetchers {
			if jwt = fetcher(r); jwt != "" {
				break
			}
		}
		if jwt == "" {
			//show error/ denied etc
			return noTokenHandler(rw, r)
		}

		idToken, err := jwtVerifier.Verify(r.Context(), jwt)
		if err != nil {
			// handle error
			switch err.(type) {
			case *oidc.TokenExpiredError:
				return expiredHandler(rw, r)
			}

			log.Printf("unkown error decoding jwt token :%v", err)
			return deniedHandler(rw, r)
		}

		r = r.WithContext(jwtoken.ToContext(r.Context(), jwtoken.New(idToken)))

		if postCheck(r) == false {
			//show error/ denied etc
			return deniedHandler(rw, r)
		}

		if next == nil {
			return true, nil
		}
		return next(rw, r)
	}
}
