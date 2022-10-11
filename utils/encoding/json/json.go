package json

import (
	"github.com/goccy/go-json"
	"io"
	"net/http"
	"strings"
)

type JsonEncoding struct{}

func NewJsonEncoding() *JsonEncoding {
	return &JsonEncoding{}
}

func (j *JsonEncoding) CanDecode(mimetype string) bool {
	mimetype = strings.ToLower(mimetype)
	return mimetype == `application/json` || strings.HasSuffix(mimetype, "+json")
}

func (j *JsonEncoding) Decode(body io.Reader, v any) error {
	return json.NewDecoder(body).Decode(v)
}

func (j *JsonEncoding) Encode(rw http.ResponseWriter, v any) error {
	return json.NewEncoder(rw).Encode(v)
}

func (j *JsonEncoding) Mimetype() string {
	return "application/json"
}
