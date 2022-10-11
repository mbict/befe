package encoding

import "io"

type Decoder interface {
	CanDecode(mimetype string) bool
	Decode(source io.Reader, v any) error
}
