package resp

import "fmt"

type integer struct{}

func (i integer) Encode(n int) []byte {
	message := fmt.Sprintf(":%d\r\n", n)
	return []byte(message)
}
