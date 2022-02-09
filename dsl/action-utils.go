package dsl

import (
	"net/http"
	"time"
)

func Delay(duration time.Duration) Action {
	return ActionFunc(func(_ http.ResponseWriter, _ *http.Request) {
		time.Sleep(duration)
	})
}
