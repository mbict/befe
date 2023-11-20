package oidc

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/vmihailenco/msgpack/v5"
)

type urlPathState struct {
	To    string `msgpack:"to,omitempty"`
	State []byte `msgpack:"state"`
}

func encodeUrlPathState(path string) string {
	state, _ := msgpack.Marshal(&urlPathState{
		To:    path,
		State: generateRandomState(8),
	})
	return base64.RawURLEncoding.EncodeToString(state)
}

func decodeUrlPathState(data string) string {
	msgPackData, err := base64.RawURLEncoding.DecodeString(data)
	if err != nil {
		return "/"
	}
	urlSt := urlPathState{}
	if err = msgpack.Unmarshal(msgPackData, &urlSt); err != nil {
		return ""
	}
	return urlSt.To
}

func generateRandomState(length int) []byte {
	token := make([]byte, length)
	rand.Read(token)
	return token
}
