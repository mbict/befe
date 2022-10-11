package xml

import (
	"encoding/xml"
	"io"
	"net/http"
	"strings"
)

type XMLEncoding struct{}

func NewXMLEncoding() *XMLEncoding {
	return &XMLEncoding{}
}

func (j *XMLEncoding) CanDecode(mimetype string) bool {
	mimetype = strings.ToLower(mimetype)
	return mimetype == `application/xml` || strings.HasSuffix(mimetype, "+xml")
}

func (j *XMLEncoding) Decode(body io.Reader, v any) error {
	return xml.NewDecoder(body).Decode(v)
}

func (j *XMLEncoding) Encode(rw http.ResponseWriter, v any) error {
	return xml.NewEncoder(rw).Encode(v)
}

func (j *XMLEncoding) Mimetype() string {
	return "application/xml"
}
