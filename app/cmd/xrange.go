package cmd

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/types"
)

func Xrange(conn net.Conn, server *types.ServerState, args ...string) {
	// Validate the number of arguments
	if len(args) < 3 {
		fmt.Printf("Expected atleast 3 arguments for 'XRANGE' command, got %v\n", args)
		return
	}

	streamKey := args[0]
	stream, ok := server.Streams[streamKey]
	if !ok {
		_, err := conn.Write(respHandler.Array.Encode([]string{}))
		if err != nil {
			fmt.Printf("Error writing response: %s\n", err)
		}
		return
	}

	items := fetchFromStream(stream, args[1], args[2])
	encodedItems, err := EncodeStreamEntrySlice(items)
	if err != nil {
		fmt.Printf("Error encoding stream entries: %s\n", err)
		return
	}
	_, err = conn.Write(encodedItems)
	if err != nil {
		fmt.Printf("Error writing response: %s\n", err)
	}
}

func fetchFromStream(streams []types.StreamEntry, start string, end string) []types.StreamEntry {
	res := make([]types.StreamEntry, 0)
	for _, stream := range streams {
		streamIDParts := strings.Split(stream.ID, "-")

		streamTimestamp, err := strconv.Atoi(streamIDParts[0])
		if err != nil {
			panic(fmt.Sprintf("Invalid ID format present: %s", stream.ID))
		}
		streamSeqNo, err := strconv.Atoi(streamIDParts[1])
		if err != nil {
			panic(fmt.Sprintf("Invalid ID format present: %s", stream.ID))
		}

		if isGreaterOrEqual(streamTimestamp, streamSeqNo, start) && isLessThanOrEqual(streamTimestamp, streamSeqNo, end) {
			res = append(res, stream)
		}
	}

	return res
}

// Helper function to check if a stream ID is greater than or equal to provided start ID
func isGreaterOrEqual(streamTimestamp int, streamSeqNo int, startID string) bool {
	if startID == "-" {
		return true
	}

	idParts := strings.Split(startID, "-")
	idTimestamp, err := strconv.Atoi(idParts[0])
	if err != nil {
		panic(fmt.Sprintf("Invalid ID format present: %s", startID))
	}

	// Handle the case where the ID is just a timestamp
	if len(idParts) == 1 {
		return streamTimestamp >= idTimestamp
	}

	idSeqNo, err := strconv.Atoi(idParts[1])
	if err != nil {
		panic(fmt.Sprintf("Invalid ID format present: %s", startID))
	}

	return streamTimestamp > idTimestamp || (streamTimestamp == idTimestamp && streamSeqNo >= idSeqNo)
}

// Helper function to check if a stream ID is less than or equal to provided end ID
func isLessThanOrEqual(streamTimestamp int, streamSeqNo int, endID string) bool {
	if endID == "+" {
		return true
	}

	idParts := strings.Split(endID, "-")
	idTimestamp, err := strconv.Atoi(idParts[0])
	if err != nil {
		panic(fmt.Sprintf("Invalid ID format present: %s", endID))
	}

	// Handle the case where the ID is just a timestamp
	if len(idParts) == 1 {
		return streamTimestamp <= idTimestamp
	}

	idSeqNo, err := strconv.Atoi(idParts[1])
	if err != nil {
		panic(fmt.Sprintf("Invalid ID format present: %s", endID))
	}

	return streamTimestamp < idTimestamp || (streamTimestamp == idTimestamp && streamSeqNo <= idSeqNo)
}
