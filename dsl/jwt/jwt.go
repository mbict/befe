package jwt

import (
	"context"
	. "github.com/mbict/befe/dsl"
	. "github.com/mbict/befe/expr"
	"github.com/mbict/befe/utils/token"
	"log"
	"net/http"
)

var jwtContextKey int

var emptyHandler = func(_ http.ResponseWriter, _ *http.Request) {}

var stopHandler = func(rw http.ResponseWriter, r *http.Request) (bool, error) {
	return false, nil
}

type Jwk interface {
	Action

	WithAudience(audiences ...string) Jwk
	WithExpiredCheck() Jwk
	WithClaim(name string, values ...string) Jwk

	WhenExpired(actions ...Action) Jwk
	WhenDenied(actions ...Action) Jwk
	WhenNoToken(actions ...Action) Jwk
}

func JwkToken(issuer string) Jwk {
	return &jwtValidator{
		issuers:           []string{issuer},
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
	issuers           []string
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

func (j *jwtValidator) BuildHandler(ctx context.Context, next Handler) Handler {
	jwtVerifier := token.NewJwkVerifier(ctx, j.issuers[0])

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
	postCheck := conditions.BuildCondition()

	return func(rw http.ResponseWriter, r *http.Request) (bool, error) {
		jwt := token.FromRequest(r)
		if jwt == nil || len(jwt) == 0 {
			//show error/ denied etc
			return noTokenHandler(rw, r)
		}

		valid, t, err := jwtVerifier.Verify(jwt)

		//default cases on errors
		if err != nil {
			switch err.Error() {
			case "exp not satisfied":
				return expiredHandler(rw, r)
			default:
				log.Printf("unkown error decoding jwt token :%v", err)
				return deniedHandler(rw, r)
			}
			return false, err
		}

		r = r.WithContext(context.WithValue(r.Context(), &jwtContextKey, t))

		if valid == false || postCheck(r) == false {
			//show error/ denied etc
			return deniedHandler(rw, r)
		}

		if next == nil {
			return true, nil
		}
		return next(rw, r)
	}
}
