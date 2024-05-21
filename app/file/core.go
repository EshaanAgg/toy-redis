package file

import (
	"fmt"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/types"
)

// Reads the length from the byte slice
// Returns the value of the length, is the value of a special type, the remaining data to process and the error if any
func readLength(data []byte) (int, bool, []byte, error) {
	if len(data) == 0 {
		return 0, false, data, fmt.Errorf("no data to read")
	}

	firstByte := data[0]
	msb := firstByte >> 6

	switch msb {
	case 0:
		// Length is present in the first byte itself
		return readIgnoringTwoMSB(firstByte), false, data[1:], nil

	case 1:
		// Length is present current 6 bits + next byte
		if len(data) < 2 {
			return 0, false, data, fmt.Errorf("not enough data to read: expected atleast 2 bytes, got %d", len(data))
		}
		return readIntIgnoringTwoMSB(data[:2]), false, data[2:], nil

	case 2:
		// Length is present in the next 4 bytes
		if len(data) < 5 {
			return 0, false, data, fmt.Errorf("not enough data to read: expected atleast 5 bytes, got %d", len(data))
		}
		return int(readAsInteger(data[1:5])), false, data[5:], nil

	case 3:
		lsb := firstByte & 0b00111111
		var bytesToRead []byte
		var nextBytes []byte

		switch lsb {
		case 0:
			// 8 bit integer follows
			if len(data) < 2 {
				return 0, false, data, fmt.Errorf("reading string encoded length: expected atleast 2 bytes, got %d", len(data))
			}
			bytesToRead = data[1:2]
			nextBytes = data[2:]

		case 1:
			// 16 bit integer
			if len(data) < 3 {
				return 0, false, data, fmt.Errorf("reading string encoded length: expected atleast 3 bytes, got %d", len(data))
			}
			bytesToRead = data[1:3]
			nextBytes = data[3:]

		case 2:
			// 32 bit integer
			if len(data) < 5 {
				return 0, false, data, fmt.Errorf("reading string encoded length: expected atleast 5 bytes, got %d", len(data))
			}
			bytesToRead = data[1:5]
			nextBytes = data[5:]

		default:
			return 0, false, data, fmt.Errorf("unimplemented value for the LSB bits %d", lsb)
		}

		return int(readAsInteger(bytesToRead)), true, nextBytes, nil

	default:
		return 0, false, data, fmt.Errorf("invalid value for the MSB: %d", msb)
	}
}

func readInteger(data []byte) (int, []byte, error) {
	n, _, remainingData, err := readLength(data)
	if err != nil {
		return 0, data, fmt.Errorf("error reading integer: %w", err)
	}
	return n, remainingData, nil
}

// Reads a string from the byte slice
// Returns the string, the remaining data to process and the error if any
func readString(data []byte) (string, []byte, error) {
	n, special, remainingData, err := readLength(data)
	if err != nil {
		return "", data, fmt.Errorf("error reading string's length: %w", err)
	}

	if special {
		return fmt.Sprintf("%d", n), remainingData, nil
	}

	if len(remainingData) < n {
		return "", data, fmt.Errorf("not enough data to read: expected atleast %d bytes, got %d", n, len(remainingData))
	}
	return string(remainingData[:n]), remainingData[n:], nil
}

// Reads the expiry for the key-value pair from the byte slice
// Returns the expiry timestamp, the remaining data to process and the error if any
// Returns -1 if no expiry timestamp is present
func readExpiry(data []byte) (int64, []byte, error) {
	if len(data) == 0 {
		return 0, data, fmt.Errorf("unexpected end of data when parsing for (optional) expiry")
	}

	switch data[0] {
	case 0xFD:
		// Expiry timestamp is present in seconds
		if len(data) < 5 {
			return 0, data, fmt.Errorf("not enough data to read: expected atleast 5 bytes, got %d", len(data))
		}
		exp := readAsInteger(data[1:5])
		return exp * 1000, data[5:], nil

	case 0xFC:
		// Expiry timestamp is present in milliseconds
		if len(data) < 9 {
			return 0, data, fmt.Errorf("not enough data to read: expected atleast 9 bytes, got %d", len(data))
		}
		return readAsInteger(data[1:9]), data[9:], nil

	default:
		// No expiry timestamp present
		return -1, data, nil
	}
}

// Reads a key-value pair from the byte slice
// Returns the key-value pair, the remaining data to process and the error if any
func readKeyValuePair(data []byte) (string, types.DBItem, []byte, error) {
	dbItem := types.DBItem{}

	if len(data) == 0 {
		return "", dbItem, data, fmt.Errorf("unexpected end of data while reading key value pairs")
	}

	expiry, remainingData, err := readExpiry(data)
	if err != nil {
		return "", dbItem, data, fmt.Errorf("error reading expiry: %w", err)
	}
	dbItem.Expiry = expiry

	if len(remainingData) == 0 {
		return "", dbItem, remainingData, fmt.Errorf("no data to read value type while reading key-value pair")
	}
	if remainingData[0] != 0 {
		return "", dbItem, remainingData, fmt.Errorf("unsupported value type provided: %d, expected 0 for string", remainingData[0])
	}

	key, remainingData, err := readString(remainingData[1:])
	if err != nil {
		return "", dbItem, data, fmt.Errorf("error reading key: %w", err)
	}
	value, remainingData, err := readString(remainingData)
	if err != nil {
		return "", dbItem, data, fmt.Errorf("error reading value: %w", err)
	}
	dbItem.Value = value

	return key, dbItem, remainingData, nil
}

func readAuxliaryField(data []byte) (string, string, []byte, error) {
	if len(data) == 0 {
		return "", "", data, fmt.Errorf("unexpected end of data while reading the auxiliary fields")
	}

	key, remainingData, err := readString(data)
	if err != nil {
		return "", "", data, fmt.Errorf("error reading key: %w", err)
	}

	value, remainingData, err := readString(remainingData)
	if err != nil {
		return "", "", data, fmt.Errorf("error reading value: %w", err)
	}

	return key, value, remainingData, nil
}

func checkHeader(data []byte) ([]byte, error) {
	if len(data) < 9 {
		return data, fmt.Errorf("not enough data to read: expected atleast 9 bytes, got %d", len(data))
	}

	if string(data[:5]) != "REDIS" {
		return data, fmt.Errorf("invalid header provided: %s", string(data[:5]))
	}

	version := string(data[5:9])
	fmt.Printf("REDIS file version: %s\n", version)
	return data[9:], nil
}

func parseFile(data []byte, serverState *types.ServerState) error {
	// Check the header of the file
	data, err := checkHeader(data)
	if err != nil {
		return fmt.Errorf("error parsing header of the file: %v", err)
	}
	if len(data) == 0 {
		return fmt.Errorf("no data to read after header in the file")
	}

	// Check the auxliary fields
	for data[0] == 0xFA {
		key, value, remainData, err := readAuxliaryField(data[1:])
		if err != nil {
			return fmt.Errorf("error reading auxliary field: %v", err)
		}
		fmt.Printf("Auxliary Field: %s -> %s\n", key, value)
		data = remainData
	}

	// Database selector
	if len(data) == 0 {
		return fmt.Errorf("file data ended abruptly after reading auxiliary fields")
	}
	if data[0] == 0xFE {
		if len(data) < 2 {
			return fmt.Errorf("not enough data to read: expected atleast 2 bytes, got %d", len(data))
		}
		fmt.Printf("Database Selector: %d\n", data[1])
		data = data[2:]
	}

	// Read the resize DB fields
	if len(data) == 0 {
		return fmt.Errorf("file data ended abruptly after reading database selector")
	}

	if data[0] == 0xFB {
		hashSize, remainData, err := readInteger(data[1:])
		if err != nil {
			return fmt.Errorf("error reading hash size: %v", err)
		}
		fmt.Printf("Hash Size: %d\n", hashSize)
		expireHashSize, remainData, err := readInteger(remainData)
		if err != nil {
			return fmt.Errorf("error reading expire hash size: %v", err)
		}
		fmt.Printf("Expire Hash Size: %d\n", expireHashSize)
		data = remainData
	}

	// Read the key-value pairs
	for len(data) > 0 && data[0] != 0xFF {
		key, dbItem, remainData, err := readKeyValuePair(data)
		if err != nil {
			return fmt.Errorf("error reading key-value pair: %v", err)
		}
		fmt.Printf("Key: %s, Value: %s, Expiry: %d\n", key, dbItem.Value, dbItem.Expiry)

		// Store the key-value pair in the database if it is not expired
		if dbItem.Expiry == -1 || dbItem.Expiry > time.Now().UnixMilli() {
			serverState.DBMutex.Lock()
			serverState.DB[key] = dbItem
			serverState.DBMutex.Unlock()
		}

		data = remainData
	}

	// TODO: Handle file checksum
	return nil
}
