package buildin

import "reflect"

//go:generate extract github.com/mbict/befe/dsl
//go:generate extract github.com/mbict/befe/utils/token

var Symbols = map[string]map[string]reflect.Value{}

func init() {
}
