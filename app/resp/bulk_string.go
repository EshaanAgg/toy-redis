package resp

import "fmt"

type bulkString struct{}

func (bulkString) Encode(s string) []byte {
	lenStr := fmt.Sprintf("%d", len(s))
	return []byte("$" + lenStr + "\r\n" + s + "\r\n")
}

// Reads the byte slice to decode the bulk string
// and returns the string and the remaining byte slice
func (bulkString) Decode(b []byte) (string, []byte, error) {
	if b[0] != '$' {
		return "", b, fmt.Errorf("expected first character as '$', got %q", b[0])
	}

	n, b, err := parseLen(b[1:])
	if err != nil {
		return "", b, fmt.Errorf("invalid format for bulk string: %v", err)
	}

	b, err = parseCRLF(b)
	if err != nil {
		return "", b, fmt.Errorf("invalid format for bulk string: %v", err)
	}

	if len(b) < n+2 {
		return "", b, fmt.Errorf("invalid format for bulk string: expected length of string to be atleast %d, got %d", n+2, len(b))
	}

	str := string(b[:n])
	b = b[n:]

	b, err = parseCRLF(b)
	if err != nil {
		return "", b, fmt.Errorf("invalid format for bulk string: %v", err)
	}

	return str, b, nil
}
