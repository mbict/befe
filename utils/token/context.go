package token

import (
	"net/http"
	"strings"
)

func FromRequest(r *http.Request) []byte {
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	if len(splitToken) != 2 {
		return nil
	}
	return []byte(splitToken[1])
}
