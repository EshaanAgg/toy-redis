package cmd

func Echo(message string) []byte {
	return respHandler.BulkStr.Encode(message)
}
