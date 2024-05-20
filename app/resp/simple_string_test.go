package resp_test

import (
	"testing"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/stretchr/testify/assert"
)

func TestStrEncode(t *testing.T) {
	var testcases = []struct {
		name        string
		input       string
		expected    []byte
		shouldError bool
	}{
		{"Simple string", "PONG", []byte("+PONG\r\n"), false},
		{"Simple string with spaces", "PONG PONG", []byte("+PONG PONG\r\n"), false},
		{"Simple string with special characters", "PONG!@#$%^&*()", []byte("+PONG!@#$%^&*()\r\n"), false},
		{"Simple string with new line", "PONG\n", nil, true},
		{"Simple string with carriage return", "PONG\r", nil, true},
		{"Simple string with carriage return and new line", "PONG\r\n", nil, true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			respHandler := resp.RESPHandler{}
			res, err := respHandler.Str.Encode(tc.input)

			if tc.shouldError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.expected, res)
			}
		})
	}
}

func TestStrDecode(t *testing.T) {
	var testcases = []struct {
		name        string
		input       []byte
		expected    string
		shouldError bool
	}{
		{"Simple string", []byte("+PONG\r\n"), "PONG", false},
		{"Simple string with spaces", []byte("+PONG PONG\r\n"), "PONG PONG", false},
		{"Simple string with special characters", []byte("+PONG!@#$%^&*()\r\n"), "PONG!@#$%^&*()", false},
		{"Simple string with new line", []byte("+PONG\n"), "", true},
		{"Simple string with carriage return", []byte("+PONG\r"), "", true},
		{"Invalid format for simple string", []byte("PONG\r\n"), "", true},
		{"Invalid format for simple string", []byte("+PONG"), "", true},
		{"Invalid format for simple string", []byte("+PONG"), "", true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			respHandler := resp.RESPHandler{}
			res, remain, err := respHandler.Str.Decode(tc.input)

			if tc.shouldError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.expected, res)
				assert.Equal(t, len(remain), 0)
			}
		})
	}
}
