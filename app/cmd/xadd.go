package cmd

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/types"
)

func Xadd(conn net.Conn, server *types.ServerState, args ...string) {
	if len(args) < 4 {
		fmt.Printf("Expected at least 4 arguments for 'XADD' command, got %v\n", args)
		return
	}

	streamKey := args[0]
	if checkIfExistsAsKV(streamKey, server) {
		fmt.Printf("Key %s already exists for a key-value pair\n", streamKey)
		return
	}

	itemKey := args[1]
	kvMap := make(map[string]string)
	for i := 2; i < len(args); i += 2 {
		kvMap[args[i]] = args[i+1]
	}

	// If the stream does not exist, create it
	if _, ok := server.Streams[streamKey]; !ok {
		fmt.Printf("Initializing stream with key: %s\n", streamKey)
		server.Streams[streamKey] = []types.StreamEntry{}
	}

	// Validate the entry ID
	var validatedEntryId string
	var validatedEntryIdErr string
	if len(server.Streams[streamKey]) > 0 {
		validatedEntryId, validatedEntryIdErr = getValidatedEntryID(
			itemKey,
			&server.Streams[streamKey][len(server.Streams[streamKey])-1].ID,
		)
	} else {
		validatedEntryId, validatedEntryIdErr = getValidatedEntryID(itemKey, nil)
	}

	// If there is an error, return the error to the client
	if validatedEntryIdErr != "" {
		fmt.Printf("Error validating entry ID: %s\n", validatedEntryIdErr)
		errBytes := respHandler.Err.Encode(validatedEntryIdErr)
		_, err := conn.Write(errBytes)
		if err != nil {
			fmt.Printf("Error writing response: %s\n", err)
		}
		return
	}

	// Add the entry to the stream
	server.Streams[streamKey] = append(server.Streams[streamKey], types.StreamEntry{
		ID:  validatedEntryId,
		KVs: kvMap,
	})

	// Return the ID of the added item
	res, err := respHandler.Str.Encode(itemKey)
	if err != nil {
		fmt.Printf("Error encoding response: %s\n", err)
		return
	}
	_, err = conn.Write(res)
	if err != nil {
		fmt.Printf("Error writing response: %s\n", err)
	}
}

// getValidatedEntryID parses the provided entry ID and returns the validated entry id and an error string if any
func getValidatedEntryID(entryID string, lastEntryID *string) (string, string) {
	entryParts := strings.Split(entryID, "-")
	if len(entryParts) != 2 {
		return "", "ERR The provided entry should have atleast two parts separated by a -"
	}

	// Parse the timestamp and sequence number
	timeStamp, err := strconv.Atoi(entryParts[0])
	if err != nil {
		return "", fmt.Sprintf("ERR Can't parse valid timestamp from entry ID: %v", err)
	}
	seqNumber, err := strconv.Atoi(entryParts[1])
	if err != nil {
		return "", fmt.Sprintf("ERR Can't parse valid sequence number from entry ID: %v", err)
	}

	// Minimum value of the ID check
	if timeStamp == 0 && seqNumber == 0 {
		return "", "ERR The ID specified in XADD must be greater than 0-0"
	}

	// If the last entry is nill, then the entry is valid
	if lastEntryID == nil {
		return entryID, ""
	}

	// Parse the last entry ID
	lastEntryParts := strings.Split(*lastEntryID, "-")
	lastTimeStamp, err := strconv.Atoi(lastEntryParts[0])
	if err != nil {
		panic("The last entry ID is not in the correct format")
	}
	lastSeqNumber, err := strconv.Atoi(lastEntryParts[1])
	if err != nil {
		panic("The last entry ID is not in the correct format")
	}

	// Check that the entry should be greater than the last entry
	if timeStamp > lastTimeStamp || (timeStamp == lastTimeStamp && seqNumber > lastSeqNumber) {
		return entryID, ""
	}

	return "", "ERR The ID specified in XADD is equal or smaller than the target stream top item"
}
