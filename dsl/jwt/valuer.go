package jwt

import (
	"github.com/mbict/befe/dsl/jwt/jwtoken"
	. "github.com/mbict/befe/expr"
	"net/http"
)

func ValueFromClaim(name string) Valuer {
	return func(r *http.Request) interface{} {
		t, ok := jwtoken.FromContext(r.Context())
		if !ok || t == nil {
			return nil
		}

		v, ok := t.Get(name)
		if !ok {
			return nil
		}

		/*
			//string value
			if value, ok := v.(string); ok {
				return value
			}

			//slice value
			if values, ok := v.([]string); ok {
				return values
			}
		*/

		return v
	}
}
