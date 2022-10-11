package dsl

import (
	"fmt"
	. "github.com/mbict/befe/expr"
	"net/http"
)

func Debug(message string, params ...interface{}) Action {
	return ActionFunc(func(rw http.ResponseWriter, r *http.Request) (bool, error) {
		fmt.Println(append([]interface{}{"[DEBUG]", message}, params...)...)
		return true, nil
	})
}
