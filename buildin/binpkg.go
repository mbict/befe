package buildin

import "reflect"

//go:generate yaegi extract github.com/mbict/befe/dsl
//go:generate yaegi extract github.com/mbict/befe/dsl/http
//go:generate yaegi extract github.com/mbict/befe/dsl/jwt
//go:generate yaegi extract github.com/mbict/befe/dsl/oidc
//go:generate yaegi extract github.com/mbict/befe/dsl/templates
//go:generate yaegi extract github.com/mbict/befe/expr

var Symbols = map[string]map[string]reflect.Value{}

func init() {
}
