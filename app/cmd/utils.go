package cmd

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/app/types"
)

func checkIfKeyExists(key string, server *types.ServerState) bool {
	_, okString := server.DB[key]
	_, okStream := server.Streams[key]

	return okString || okStream
}

func checkIfExistsAsKV(key string, server *types.ServerState) bool {
	_, ok := server.DB[key]
	return ok
}

func EncodeStreamEntrySlice(entries []types.StreamEntry) ([]byte, error) {
	encodedBytes := []byte(fmt.Sprintf("*%d\r\n", len(entries)))

	for _, entry := range entries {
		encodedEntry, err := EncodeStreamEntry(entry)
		if err != nil {
			return nil, fmt.Errorf("failed to encode entry: %v", err)
		}
		encodedBytes = append(encodedBytes, encodedEntry...)
	}

	return encodedBytes, nil
}

func EncodeStreamEntry(entry types.StreamEntry) ([]byte, error) {
	encodedBytes := []byte(fmt.Sprintf("*%d\r\n", 2))

	// Add the ID to the encoded bytes
	encodedID, err := respHandler.Str.Encode(entry.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to encode ID: %v", err)
	}
	encodedBytes = append(encodedBytes, encodedID...)

	// Add the key-value pairs to the encoded bytes
	kvSlice := make([]string, 0)
	for key, value := range entry.KVs {
		kvSlice = append(kvSlice, key, value)
	}
	encodedKVs := respHandler.Array.Encode(kvSlice)
	encodedBytes = append(encodedBytes, encodedKVs...)

	return encodedBytes, nil
}

func EncodeStreamResult(result StreamResult) ([]byte, error) {
	encodedBytes := []byte(fmt.Sprintf("*%d\r\n", 2))

	// Add the stream key to the encoded bytes
	encodedStreamKey, err := respHandler.Str.Encode(result.StreamKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encode stream key: %v", err)
	}
	encodedBytes = append(encodedBytes, encodedStreamKey...)

	// Add the stream entries to the encoded bytes
	encodedEntries, err := EncodeStreamEntrySlice(result.StreamEntries)
	if err != nil {
		return nil, fmt.Errorf("failed to encode stream entries: %v", err)
	}
	encodedBytes = append(encodedBytes, encodedEntries...)

	return encodedBytes, nil
}

func EncodeStreamResultArray(results []StreamResult) ([]byte, error) {
	encodedBytes := []byte(fmt.Sprintf("*%d\r\n", len(results)))

	for _, result := range results {
		encodedResult, err := EncodeStreamResult(result)
		if err != nil {
			return nil, fmt.Errorf("failed to encode result: %v", err)
		}
		encodedBytes = append(encodedBytes, encodedResult...)
	}

	return encodedBytes, nil
}
