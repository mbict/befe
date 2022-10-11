package encoding

type ContentDecoder []Decoder

func (d *ContentDecoder) GetDecoder(mimetype string) Decoder {
	for _, dec := range *d {
		if dec.CanDecode(mimetype) == true {
			return dec
		}
	}
	return nil
}
