package cmd

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/types"
)

func Xadd(conn net.Conn, server *types.ServerState, args ...string) {
	if len(args) < 4 {
		fmt.Printf("Expected at least 4 arguments for 'XADD' command, got %v\n", args)
		return
	}

	streamKey := args[0]
	fmt.Println("streamKey: ", streamKey)
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
	server.Streams[streamKey] = append(server.Streams[streamKey], types.StreamEntry{
		ID:  itemKey,
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
