package resp

import "fmt"

type simpleString struct{}

func (simpleString) Encode(s string) ([]byte, error) {
	// The string must not contain \r or \n
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' || s[i] == '\r' {
			return nil, fmt.Errorf("invalid character in simple string: %c", s[i])
		}
	}

	return []byte("+" + s + "\r\n"), nil
}

func (simpleString) Decode(b []byte) (string, error) {
	if b[0] != '+' {
		return "", fmt.Errorf("invalid format for simple string: Expected the first byte to be '+', got '%c'", b[0])
	}

	err := checkLastTwoBytes(b)
	if err != nil {
		return "", err
	}

	for i := 1; i < len(b)-2; i++ {
		if b[i] == '\r' || b[i] == '\n' {
			return "", fmt.Errorf("invalid format for simple string: The string content contains %c, which is not allowed", b[i])
		}
	}

	return string(b[1 : len(b)-2]), nil
}
