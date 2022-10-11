package expr

import (
	"github.com/mbict/befe/utils/encoding"
	"github.com/mbict/befe/utils/encoding/json"
	"github.com/mbict/befe/utils/encoding/xml"
	"io"
	"net/http"
)

var RequestDecoder = encoding.ContentDecoder{
	json.NewJsonEncoding(),
	xml.NewXMLEncoding(),
}

func DecodeResponse(response *http.Response) (interface{}, error) {
	if decoder := RequestDecoder.GetDecoder(response.Header.Get("Content-Type")); decoder != nil {
		var data interface{}
		err := decoder.Decode(response.Body, &data)
		defer response.Body.Close()
		return data, err
	} else {
		data, err := io.ReadAll(response.Body)
		defer response.Body.Close()
		return data, err
	}
}
