package resp

import (
	"fmt"
	"strconv"
)

type integer struct{}

func (integer) Encode(n int) []byte {
	message := fmt.Sprintf(":%d\r\n", n)
	return []byte(message)
}

func (integer) Decode(b []byte) (int, error) {
	if b[0] != ':' {
		return 0, fmt.Errorf("not correct encoding for int, should begin with :, but does with %b", b[0])
	}

	l := len(b)
	if string(b[l-2:]) != "\r\n" {
		return 0, fmt.Errorf("error value must end in CRLF, but it ends in %s", string(b[l-2:]))
	}

	n, err := strconv.Atoi(string(b[1 : l-2]))
	if err != nil {
		return 0, fmt.Errorf("can't convert the passed bytes to integer: %v", err)
	}

	return n, err
}
