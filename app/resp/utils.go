package resp

import (
	"fmt"
	"strconv"
)

// checkLastTwoBytes checks if the last two bytes of the byte slice are '\r\n'
func checkLastTwoBytes(b []byte) error {
	if b[len(b)-2] == '\r' && b[len(b)-1] == '\n' {
		return nil
	}
	return fmt.Errorf("invalid format for simple string: Expected the last two bytes to be \\r\\n, but got %c%c", b[len(b)-2], b[len(b)-1])
}

// Parses the length of the string from the byte slice and returns the length and the remaining byte slice
func parseLen(b []byte) (int, []byte, error) {
	lenStr := ""
	for i := 0; i < len(b); i++ {
		if b[i] == '\r' {
			break
		}
		lenStr += string(b[i])
	}

	n, err := strconv.Atoi(lenStr)
	if err != nil {
		return 0, nil, fmt.Errorf("cannot parse number from string %s: %v", lenStr, err)
	}

	return n, b[len(lenStr):], nil
}

// Checks if the first two bytes of the byte slice are '\r\n'
// and returns the remaining byte slice
func parseCRLF(b []byte) ([]byte, error) {
	if b[0] != '\r' || b[1] != '\n' {
		return nil, fmt.Errorf("expected the next two bytes to be \\r\\n, got %c%c", b[0], b[1])
	}

	return b[2:], nil
}
