package dsl

import (
	"context"
	"github.com/mbict/befe/utils/token"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

var jwtContextKey int

type Jwk interface {
	//Condition
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
		issuers: []string{issuer},
	}
}

type jwtValidator struct {
	issuers           []string
	withAudienceCheck []string
	withExpiredCheck  bool
	withClaims        map[string][]string

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
	deniedHandler := j.deniedActions.BuildHandler(ctx, emptyHandler)

	if j.expiredActions == nil {
		j.expiredActions = Actions{Deny()}
	}
	expiredHandler := j.expiredActions.BuildHandler(ctx, emptyHandler)

	if j.noTokenActions == nil {
		j.noTokenActions = Actions{Unauthorized()}
	}
	noTokenHandler := j.noTokenActions.BuildHandler(ctx, emptyHandler)

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

	return func(rw http.ResponseWriter, r *http.Request) {
		jwt := token.FromRequest(r)
		if jwt == nil || len(jwt) == 0 {
			//show error/ denied etc
			noTokenHandler(rw, r)
			return
		}

		valid, t, err := jwtVerifier.Verify(jwt)

		//default cases on errors
		if err != nil {
			switch err.Error() {
			case "exp not satisfied":
				expiredHandler(rw, r)
			default:
				log.Printf("unkown error decoding jwt token :%v", err)
				deniedHandler(rw, r)
			}
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), &jwtContextKey, t))

		if valid == false || postCheck(r) == false {
			//show error/ denied etc
			deniedHandler(rw, r)
			return
		}

		next(rw, r)
	}
}

//--- conditional dsl
func HasJwtClaim(name string, values ...string) Condition {
	return ConditionHandler(func(r *http.Request) bool {
		t, ok := r.Context().Value(&jwtContextKey).(*token.JwtToken)
		if !ok {
			return false
		}

		v, ok := t.Get(name)
		if !ok {
			return false
		}

		//string value
		if value, ok := v.(string); ok {
			for _, valueMatch := range values {
				if valueMatch == value {
					return true
				}
			}
		}

		//slice value
		if claimValues, ok := v.([]string); ok {
			for _, value := range claimValues {
				for _, valueMatch := range values {
					if valueMatch == value {
						return true
					}
				}
			}
		}
		return false
	})
}

func HasJwtScopes(scopes ...string) Condition {
	hasScope := func(scopes []string, scope string) bool {
		for _, matchScope := range scopes {
			if scope == matchScope {
				return true
			}
		}
		return false
	}

	hasScopeInterface := func(scopes []interface{}, scope string) bool {
		for _, matchScope := range scopes {
			if scope == matchScope {
				return true
			}
		}
		return false
	}

	return ConditionHandler(func(r *http.Request) bool {
		t, ok := r.Context().Value(&jwtContextKey).(*token.JwtToken)
		if !ok {
			return false
		}

		v, ok := t.Get("scp")
		if !ok {
			return false
		}

		//string value
		for _, matchScope := range scopes {
			switch claimValue := v.(type) {
			case string:
				claimScopes := strings.Split(claimValue, " ")

				if hasScope(claimScopes, matchScope) == false {
					return false
				}

				return true

			case []string:
				if hasScope(claimValue, matchScope) == false {
					return false
				}
				return true

			case []interface{}:
				if hasScopeInterface(claimValue, matchScope) == false {
					return false
				}
				return true

			default:
				panic("unkown jwt type to handle")
			}
		}
		return false
	})
}

//HasJwtAudience matches if any of the audiences supplied match one of the audiences in the token
// Global matching is possible matching api://test/* will match api://test/1234
func HasJwtAudience(audiences ...string) Condition {
	return ConditionHandler(func(r *http.Request) bool {
		t, ok := r.Context().Value(&jwtContextKey).(*token.JwtToken)
		if !ok {
			return false
		}

		for _, audience := range t.Audience() {
			for _, audienceMatch := range audiences {
				match, err := filepath.Match(audienceMatch, audience)
				if match == true && err == nil {
					return true
				}
			}
		}
		return false
	})
}

func JwtIsNotExpired() Condition {
	return ConditionHandler(func(r *http.Request) bool {
		t, ok := r.Context().Value(&jwtContextKey).(*token.JwtToken)
		if !ok {
			return false
		}
		return t.IsExpired()
	})
}
