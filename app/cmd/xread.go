package cmd

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

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

	toBlock := false
	timeToBlock := -1

	if strings.ToUpper(args[0]) == "BLOCK" {
		time, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Printf("Failed to parse time to block: %v\n", err)
			return
		}
		args = args[2:]
		toBlock = true
		timeToBlock = time
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

	existingResult := getStreamResults(server, args)
	if toBlock {
		if timeToBlock == 0 {
			blockTillNewResult(server, args, existingResult)
		} else {
			time.Sleep(time.Duration(timeToBlock) * time.Millisecond)
		}
	}

	result := getStreamResults(server, args)
	for i := len(args) / 2; i < len(args); i++ {
		if args[i] == "$" {
			result = getNewResults(existingResult, result)
			break
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

func getNewResults(baseResult []StreamResult, newResult []StreamResult) []StreamResult {
	var result []StreamResult

	for _, newStreamResult := range newResult {
		found := false

		for _, baseStreamResult := range baseResult {
			if newStreamResult.StreamKey == baseStreamResult.StreamKey {
				found = true
				result = append(result, StreamResult{
					StreamKey:     newStreamResult.StreamKey,
					StreamEntries: getNewStreamEntries(baseStreamResult.StreamEntries, newStreamResult.StreamEntries),
				})
			}
		}

		if !found {
			result = append(result, newStreamResult)
		}
	}
	return result
}

func getNewStreamEntries(baseStreamEntries []types.StreamEntry, newStreamEntries []types.StreamEntry) []types.StreamEntry {
	var result []types.StreamEntry

	for _, newStreamEntry := range newStreamEntries {
		found := false
		for _, baseStreamEntry := range baseStreamEntries {
			if newStreamEntry.ID == baseStreamEntry.ID {
				found = true
				break
			}
		}
		if !found {
			result = append(result, newStreamEntry)
		}
	}

	return result
}

func blockTillNewResult(server *types.ServerState, args []string, baseResult []StreamResult) {
	for {
		result := getStreamResults(server, args)
		if !isSameResult(baseResult, result) {
			return
		}
		// Sleep for 20 milliseconds before checking again to avoid busy waiting
		time.Sleep(20 * time.Millisecond)
	}
}

func isSameResult(baseResult []StreamResult, newResult []StreamResult) bool {
	if len(baseResult) != len(newResult) {
		return false
	}

	for i := 0; i < len(baseResult); i++ {
		if baseResult[i].StreamKey != newResult[i].StreamKey {
			return false
		}

		if len(baseResult[i].StreamEntries) != len(newResult[i].StreamEntries) {
			return false
		}

		for j := 0; j < len(baseResult[i].StreamEntries); j++ {
			if baseResult[i].StreamEntries[j].ID != newResult[i].StreamEntries[j].ID {
				return false
			}

			if len(baseResult[i].StreamEntries[j].KVs) != len(newResult[i].StreamEntries[j].KVs) {
				return false
			}

			for key, value := range baseResult[i].StreamEntries[j].KVs {
				if newValue, ok := newResult[i].StreamEntries[j].KVs[key]; !ok || value != newValue {
					return false
				}
			}
		}
	}

	return true
}

func getStreamResults(server *types.ServerState, args []string) []StreamResult {
	numberOfStreams := len(args) / 2
	result := make([]StreamResult, numberOfStreams)

	for i := 0; i < numberOfStreams; i++ {
		streamKey := args[i]
		stream, ok := server.Streams[streamKey]

		// If the stream does not exist, then skip it
		if !ok {
			continue
		}

		startKey := args[i+numberOfStreams]

		var streamEntries []types.StreamEntry
		if startKey == "$" {
			streamEntries = stream
		} else {
			streamEntries = fetchFromStreamTillEnd(stream, startKey)
		}

		result[i] = StreamResult{
			StreamKey:     streamKey,
			StreamEntries: streamEntries,
		}
	}

	return result
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
