package resp

import "fmt"

type array struct{}

func (array) Encode(arr []string) []byte {
	bulkString := bulkString{}
	byteSlice := []byte("*" + fmt.Sprintf("%d", len(arr)) + "\r\n")

	for _, str := range arr {
		byteSlice = append(byteSlice, bulkString.Encode(str)...)
	}

	return byteSlice
}

func (array) EncodeFromElementBytes(arr [][]byte) []byte {
	byteSlice := []byte("*" + fmt.Sprintf("%d", len(arr)) + "\r\n")

	for _, bytes := range arr {
		byteSlice = append(byteSlice, bytes...)
	}

	return byteSlice
}

func (array) Decode(b []byte) ([]string, []byte, error) {
	if b[0] != '*' {
		return nil, b, fmt.Errorf("invalid format for array: expected the first byte to be '*', got '%q'", b[0])
	}

	n, b, err := parseLen(b[1:])
	if err != nil {
		return nil, b, fmt.Errorf("invalid format for array: %v", err)
	}

	b, err = parseCRLF(b)
	if err != nil {
		return nil, b, fmt.Errorf("invalid format for array: %v", err)
	}

	arr := make([]string, n)
	bulkString := bulkString{}

	str := ""
	for i := 0; i < n; i++ {
		str, b, err = bulkString.Decode(b)
		if err != nil {
			return nil, b, fmt.Errorf("invalid format for array: %v", err)
		}
		arr[i] = str
	}

	return arr, b, nil
}
