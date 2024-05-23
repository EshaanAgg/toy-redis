package cmd

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/types"
)

type StreamResult struct {
	StreamKey     string
	StreamEntries []types.StreamEntry
}

func Xread(conn net.Conn, server *types.ServerState, args ...string) {
	// Validate the number of arguments
	if len(args) < 3 {
		fmt.Printf("Expected atleast 3 arguments for 'XREAD' command, got %v\n", args)
		return
	}

	if strings.ToUpper(args[0]) != "STREAMS" {
		fmt.Printf("Expected 'STREAMS' as first argument for 'XREAD' command, got %v\n", args[0])
		return
	}

	args = args[1:]
	if len(args)%2 != 0 {
		fmt.Printf("Expected even number of arguments after 'STREAMS' for 'XREAD' command, got %v\n", args[1:])
		return
	}

	numberOfStreams := len(args) / 2
	result := make([]StreamResult, numberOfStreams)

	for i := 0; i < numberOfStreams; i++ {
		streamKey := args[i]
		stream, ok := server.Streams[streamKey]

		// If the stream does not exist, add an empty result
		if !ok {
			result[i] = StreamResult{
				StreamKey:     streamKey,
				StreamEntries: []types.StreamEntry{},
			}
			continue
		}

		startKey := args[i+numberOfStreams]
		streamEntries := fetchFromStreamTillEnd(stream, startKey)
		result[i] = StreamResult{
			StreamKey:     streamKey,
			StreamEntries: streamEntries,
		}
	}

	// Encode the result and write it to the connection
	encodedResult, err := EncodeStreamResultArray(result)
	if err != nil {
		fmt.Printf("Failed to encode result: %v\n", err)
		return
	}
	_, err = conn.Write(encodedResult)
	if err != nil {
		fmt.Printf("Failed to write result to connection: %v\n", err)
	}
}

func fetchFromStreamTillEnd(streams []types.StreamEntry, start string) []types.StreamEntry {
	startIndex := len(streams) // Set the start index to the end of the stream

	for ind, stream := range streams {
		streamIDParts := strings.Split(stream.ID, "-")

		streamTimestamp, err := strconv.Atoi(streamIDParts[0])
		if err != nil {
			panic(fmt.Sprintf("Invalid ID format present: %s", stream.ID))
		}
		streamSeqNo, err := strconv.Atoi(streamIDParts[1])
		if err != nil {
			panic(fmt.Sprintf("Invalid ID format present: %s", stream.ID))
		}

		if isGreater(streamTimestamp, streamSeqNo, start) {
			startIndex = ind
			break
		}
	}

	return streams[startIndex:]
}

// Helper function to check if a stream ID is greater than provided start ID
func isGreater(streamTimestamp int, streamSeqNo int, startID string) bool {
	idParts := strings.Split(startID, "-")
	idTimestamp, err := strconv.Atoi(idParts[0])
	if err != nil {
		panic(fmt.Sprintf("Invalid ID format present: %s", startID))
	}

	var idSeqNo int
	if len(idParts) == 1 {
		// If the ID is just a timestamp, then the sequence number is 0
		idSeqNo = 0
	} else {
		idSeqNo, err = strconv.Atoi(idParts[1])
		if err != nil {
			panic(fmt.Sprintf("Invalid ID format present: %s", startID))
		}

	}

	return streamTimestamp > idTimestamp || (streamTimestamp == idTimestamp && streamSeqNo > idSeqNo)
}
