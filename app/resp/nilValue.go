package resp

type nilValue struct{}

func (n nilValue) Encode() []byte {
	return []byte("-1\r\n")
}
