package resp_test

import (
	"testing"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/stretchr/testify/assert"
)

func TestArrayEncode(t *testing.T) {
	testcases := []struct {
		name     string
		input    []string
		expected []byte
	}{
		{"Array with 1 element", []string{"A"}, []byte("*1\r\n$1\r\nA\r\n")},
		{"Array with 2 elements", []string{"A", "B"}, []byte("*2\r\n$1\r\nA\r\n$1\r\nB\r\n")},
		{"Empty array", []string{}, []byte("*0\r\n")},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			respHandler := resp.RESPHandler{}
			res := respHandler.Array.Encode(tc.input)

			assert.Equal(t, tc.expected, res)
		})
	}
}

func TestArrayDecode(t *testing.T) {
	testcases := []struct {
		name        string
		input       []byte
		expected    []string
		shouldError bool
	}{
		{"Array with 1 element", []byte("*1\r\n$1\r\nA\r\n"), []string{"A"}, false},
		{"Array with 2 elements", []byte("*2\r\n$1\r\nA\r\n$1\r\nB\r\n"), []string{"A", "B"}, false},
		{"Empty array", []byte("*0\r\n"), []string{}, false},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			respHandler := resp.RESPHandler{}
			res, rem, err := respHandler.Array.Decode(tc.input)

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
