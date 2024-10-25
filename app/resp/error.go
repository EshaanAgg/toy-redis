package resp

import "fmt"

type errorValue struct{}

func (errorValue) Encode(err string) []byte {
	return []byte("-" + err + "\r\n")
}

func (errorValue) Decode(b []byte) (string, error) {
	if b[0] != '-' {
		return "", fmt.Errorf("not a correct error value, should begin with -, but does with %b", b[0])
	}

	l := len(b)
	if string(b[l-2:]) != "\r\n" {
		return "", fmt.Errorf("error value must end in CRLF, but it ends in %s", string(b[l-2:]))
	}

	return string(b[1 : l-2]), nil
}
