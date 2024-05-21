package file

// Used to read the number represented by a byte slice
func readAsInteger(data []byte) int64 {
	res := int64(0)

	for i := len(data) - 1; i >= 0; i-- {
		res = res<<8 + int64(data[i])
	}

	return res
}

// Used to read the number represented by the 6 LSB of a byte
func readIgnoringTwoMSB(b byte) int {
	mask := byte(0x3F) // 0011 1111
	return int(b & mask)
}

// Used to read the number represented by a byte slice, ignoring the 2 MSB of the first byte
func readIntIgnoringTwoMSB(data []byte) int {
	res := 0

	for i := len(data) - 1; i >= 0; i-- {
		res = res << 8
		if i == 0 {
			res += readIgnoringTwoMSB(data[i])
		} else {
			res += int(data[i])
		}
	}

	return res
}
