package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/types"
)

func Xadd(server *types.ServerState, args ...string) []byte {
	if len(args) < 4 {
		return respHandler.Err.Encode(
			fmt.Sprintf("ERR wrong number of arguments for '%s' command, expected at least 4, for %d", args[0], len(args)),
		)
	}

	streamKey := args[0]
	if checkIfExistsAsKV(streamKey, server) {
		return respHandler.Err.Encode("ERR key already exists for a key-value pair")
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
		return respHandler.Err.Encode(validatedEntryIdErr)
	}

	// Add the entry to the stream
	server.Streams[streamKey] = append(server.Streams[streamKey], types.StreamEntry{
		ID:  validatedEntryId,
		KVs: kvMap,
	})

	// Return the ID of the added item
	return respHandler.BulkStr.Encode(validatedEntryId)
}

// getValidatedEntryID parses the provided entry ID and returns the validated entry id and an error string if any
func getValidatedEntryID(entryID string, lastEntryID *string) (string, string) {
	if lastEntryID == nil {
		return handleNoLastEntry(entryID)
	}

	lastEntryParts := strings.Split(*lastEntryID, "-")
	lastTimeStamp, err := strconv.Atoi(lastEntryParts[0])
	if err != nil {
		panic("The last entry ID is not in the correct format")
	}
	lastSeqNumber, err := strconv.Atoi(lastEntryParts[1])
	if err != nil {
		panic("The last entry ID is not in the correct format")
	}

	return handleWithLastEntry(entryID, int64(lastTimeStamp), lastSeqNumber)
}

// Provides the validated entry ID when there is no last entry
func handleNoLastEntry(entryID string) (string, string) {
	// Generate the new entry ID completely
	if entryID == "*" {
		return fmt.Sprintf("%d-0", time.Now().UnixMilli()), ""
	}

	entryParts := strings.Split(entryID, "-")
	if len(entryParts) != 2 {
		return "", "ERR The provided entry should have atleast two parts separated by a -"
	}

	// Use the provided timestamp and use 0 as the sequence number
	if entryParts[1] == "*" {
		var seqID int
		if entryParts[0] == "0" {
			seqID = 1
		} else {
			seqID = 0
		}
		return fmt.Sprintf("%s-%d", entryParts[0], seqID), ""
	}

	// Handle the only invalid case
	if entryParts[0] == "0" && entryParts[1] == "0" {
		return "", "ERR The ID specified in XADD must be greater than 0-0"
	}

	return entryID, ""
}

// Provides the validated entry ID when there is a last entry
func handleWithLastEntry(entryID string, lastTimestamp int64, lastSeqNo int) (string, string) {
	// Generate the new entry ID completely
	if entryID == "*" {
		currTimestamp := time.Now().UnixMilli()
		curSeqNo := 0
		if currTimestamp < lastTimestamp {
			return "", "ERR The timestamp part of the entry ID is smaller than the last entry's timestamp"
		}
		if currTimestamp == lastTimestamp {
			curSeqNo = lastSeqNo + 1
		}
		return fmt.Sprintf("%d-%d", currTimestamp, curSeqNo), ""
	}

	entryParts := strings.Split(entryID, "-")
	if len(entryParts) != 2 {
		return "", "ERR The provided entry should have atleast two parts separated by a -"
	}

	// Generate the sequence number based on the last entry
	if entryParts[1] == "*" {
		parsedTime, err := strconv.Atoi(entryParts[0])
		if err != nil {
			return "", "ERR The timestamp part of the entry ID is not a valid integer"
		}
		currTimestamp := int64(parsedTime)
		if currTimestamp < lastTimestamp {
			return "", "ERR The timestamp part of the entry ID is smaller than the last entry's timestamp"
		}
		if currTimestamp == lastTimestamp {
			return fmt.Sprintf("%d-%d", currTimestamp, lastSeqNo+1), ""
		}
		return fmt.Sprintf("%d-0", currTimestamp), ""
	}

	// Handle the case when both are provided
	entryTimestamp, err := strconv.Atoi(entryParts[0])
	if err != nil {
		return "", "ERR The timestamp part of the entry ID is not a valid integer"
	}
	entrySeqNo, err := strconv.Atoi(entryParts[1])
	if err != nil {
		return "", "ERR The sequence number part of the entry ID is not a valid integer"
	}

	if entryTimestamp == 0 && entrySeqNo == 0 {
		return "", "ERR The ID specified in XADD must be greater than 0-0"
	}
	if entryTimestamp < int(lastTimestamp) || (entryTimestamp == int(lastTimestamp) && entrySeqNo <= lastSeqNo) {
		return "", "ERR The ID specified in XADD is equal or smaller than the target stream top item"
	}

	return entryID, ""
}
