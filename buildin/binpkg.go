package buildin

import "reflect"

//go:generate extract github.com/mbict/befe/dsl
//go:generate extract github.com/mbict/befe/dsl/http
//go:generate extract github.com/mbict/befe/dsl/jwt
//go:generate extract github.com/mbict/befe/dsl/oidc
//go:generate extract github.com/mbict/befe/dsl/templates
//go:generate extract github.com/mbict/befe/expr

var Symbols = map[string]map[string]reflect.Value{}

func init() {
}
