package resp

import "fmt"

func checkLastTwoBytes(b []byte) error {
	if b[len(b)-2] == '\r' && b[len(b)-1] == '\n' {
		return nil
	}
	return fmt.Errorf("invalid format for simple string: Expected the last two bytes to be \\r\\n, but got %c%c", b[len(b)-2], b[len(b)-1])
}
