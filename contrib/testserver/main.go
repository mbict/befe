package main

import (
	"log"
	"net/http"
	"strings"
)

func main() {
	log.Printf("starting example server on port :%d", 8091)
	var bigPayload = []byte(`{ "data": [{"this":"is", "a":"fake", "json":[ "response", "from", "the", "webserver"]}` + strings.Repeat(`,{"this":"is", "a":"fake", "json":[ "response", "from", "the", "webserver"]}`+"\n", 100) + `] }`)
	http.ListenAndServe(":8091", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json; charset=utf-8")
		rw.WriteHeader(200)
		rw.Write(bigPayload)
	}))
}
