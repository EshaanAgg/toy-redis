package resp

import (
	"fmt"
	"strings"
)

type RESPHandler struct {
	Str     simpleString
	BulkStr bulkString
	Array   array
	Nil     nilValue
	Int     integer
	Err     errorValue
}

func (h RESPHandler) DecodeCommand(b []byte) ([]string, []byte, error) {
	arr, next, err := h.Array.Decode(b)
	return arr, next, err
}

func (h RESPHandler) DecodeResponse(b []byte) (string, error) {
	if len(b) == 0 {
		return "", fmt.Errorf("the provided byte array is empty")
	}

	var res string
	var err error

	switch b[0] {
	case '+':
		res, b, err = h.Str.Decode(b)

	case '*':
		resArr, leftBytes, error := h.Array.Decode(b)
		err = error
		b = leftBytes
		res = fmt.Sprintf("[%s]", strings.Join(resArr, ", "))

	case '$':
		if string(b) == NIL_VALUE {
			return "NIL", nil
		}
		res, b, err = h.BulkStr.Decode(b)

	case '-':
		res, err = h.Err.Decode(b)
		b = make([]byte, 0)

	case ':':
		n, error := h.Int.Decode(b)
		err = error
		res = fmt.Sprint(n)
		b = make([]byte, 0)

	default:
		return "", fmt.Errorf("can't recognize the encoding for %s", string(b))
	}

	if err != nil {
		return "", fmt.Errorf("there was an error in parsing the response: %v", err)
	}
	if len(b) != 0 {
		return "", fmt.Errorf("recieved initial response '%s', but then recieved additional bytes: %s", res, string(b))
	}

	return res, nil
}
