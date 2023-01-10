package jwt

import (
	"github.com/mbict/befe/dsl/jwt/jwtoken"
	. "github.com/mbict/befe/expr"
	"net/http"
	"path/filepath"
	"strings"
)

func HasJwtClaim(name string, values ...string) Condition {
	return ConditionFunc(func(r *http.Request) bool {
		t, ok := jwtoken.FromContext(r.Context())
		if !ok {
			return false
		}

		v, ok := t.Get(name)
		if !ok {
			return false
		}

		//no values, we only check if the claim exists
		if len(values) == 0 {
			return true
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

	return ConditionFunc(func(r *http.Request) bool {
		t, ok := jwtoken.FromContext(r.Context())
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

// HasJwtAudience matches if any of the audiences supplied match one of the audiences in the token
// Global matching is possible matching api://test/* will match api://test/1234
func HasJwtAudience(audiences ...string) Condition {
	return ConditionFunc(func(r *http.Request) bool {
		t, ok := jwtoken.FromContext(r.Context())
		if !ok {
			return false
		}

		for _, audience := range t.Audience {
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
	return ConditionFunc(func(r *http.Request) bool {
		t, ok := jwtoken.FromContext(r.Context())
		if !ok {
			return false
		}
		return t.IsExpired()
	})
}
