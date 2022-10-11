package expr

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Valuer func(r *http.Request) interface{}

// ValueToString helper to convert any value to a string
func ValueToString(in interface{}) string {
	if in == nil {
		return ""
	}

	switch v := in.(type) {
	case string:
		return v
	case []interface{}:
		res := make([]string, len(v))
		for i, val := range v {
			res[i] = ValueToString(val)
		}
		return ValueToString(res)
	case []string:
		return strings.Join(v, `,`)
	case int:
		return strconv.Itoa(v)
	case float64:
		return fmt.Sprintf("%g", v)
	case bool:
		if v == true {
			return "true"
		}
		return "false"
	default:
		return fmt.Sprintf("%s", v)
	}
}
