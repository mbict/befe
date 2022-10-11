package expr

import "net/http"

type Param func(r *http.Request) (string, interface{})
