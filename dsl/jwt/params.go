package jwt

import (
	. "github.com/mbict/befe/dsl"
	. "github.com/mbict/befe/expr"
)

func ParamFromClaim(paramName string, claimName string) Param {
	return WithParam(paramName, ValueFromClaim(claimName))
}
