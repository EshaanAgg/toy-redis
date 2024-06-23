package resp

type nilValue struct{}

const NIL_VALUE = "$-1\r\n"

func (n nilValue) Encode() []byte {
	return []byte(NIL_VALUE)
}
