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

func Redirect(url string, status int) Action {
	return ActionFunc(http.RedirectHandler(url, status).ServeHTTP)
}

func TemporaryRedirect(url string) Action {
	return Redirect(url, http.StatusTemporaryRedirect)
}

func PermanentRedirect(url string, status int) Action {
	return Redirect(url, http.StatusPermanentRedirect)
}
