package encoding

import (
	"net/http"
)

type Encoder interface {
	Encode(rw http.ResponseWriter, v any) error
	Mimetype() string
}
