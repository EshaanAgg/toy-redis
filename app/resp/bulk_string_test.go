package resp_test

import (
	"testing"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/stretchr/testify/assert"
)

func TestBulkStrEncode(t *testing.T) {
	testcases := []struct {
		name     string
		input    string
		expected []byte
	}{
		{"Bulk string", "PONG", []byte("$4\r\nPONG\r\n")},
		{"Bulk string with spaces", "PONG PONG", []byte("$9\r\nPONG PONG\r\n")},
		{"Bulk string with new line", "PONG\n", []byte("$5\r\nPONG\n\r\n")},
		{"Bulk string with carriage return", "PONG\r", []byte("$5\r\nPONG\r\r\n")},
		{"Bulk string with carriage return and new line", "PONG\r\n", []byte("$6\r\nPONG\r\n\r\n")},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			respHandler := resp.RESPHandler{}
			res := respHandler.BulkStr.Encode(tc.input)
			assert.Equal(t, tc.expected, res)
		})
	}
}

func TestBulkStrDecode(t *testing.T) {
	var testcases = []struct {
		name        string
		input       []byte
		expected    string
		shouldError bool
	}{
		{"Bulk string", []byte("$4\r\nPONG\r\n"), "PONG", false},
		{"Bulk string with spaces", []byte("$9\r\nPONG PONG\r\n"), "PONG PONG", false},
		{"Bulk string with new line", []byte("$5\r\nPONG\n\r\n"), "PONG\n", false},
		{"Bulk string with carriage return", []byte("$5\r\nPONG\r\r\n"), "PONG\r", false},
		{"Bulk string with carriage return and new line", []byte("$6\r\nPONG\r\n\r\n"), "PONG\r\n", false},
		{"Invalid bulk string", []byte("$4\r\nPONG\r"), "", true},
		{"Invalid bulk string with missing new line", []byte("$4\r\nPONG"), "", true},
		{"Invalid bulk string with missing carriage return", []byte("$4\nPONG\r\n"), "", true},
		{"Invalid bulk string with missing carriage return and new line", []byte("$4\nPONG"), "", true},
		{"Invalid bulk string with missing length", []byte("$\r\nPONG\r\n"), "", true},
		{"Invalid bulk string with missing length and carriage return", []byte("$\nPONG\r\r\n"), "", true},
		{"Invalid bulk string with missing length and new line", []byte("$\nPONG\n\r\n"), "", true},
		{"Empty bulk string", []byte("$0\r\n\r\n"), "", false},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			respHandler := resp.RESPHandler{}
			res, rem, err := respHandler.BulkStr.Decode(tc.input)

			if tc.shouldError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.expected, res)
				assert.Equal(t, 0, len(rem))
			}
		})
	}
}
