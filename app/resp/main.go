package resp

type RESPHandler struct {
	Str     simpleString
	BulkStr bulkString
	Array   array
}

func (h RESPHandler) DecodeCommand(b []byte) ([]string, error) {
	arr, _, err := h.Array.Decode(b)
	return arr, err
}
