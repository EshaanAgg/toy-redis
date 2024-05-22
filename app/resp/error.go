package resp

type errorValue struct{}

func (errorValue) Encode(err string) []byte {
	return []byte("-" + err + "\r\n")
}
