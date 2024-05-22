package resp

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
