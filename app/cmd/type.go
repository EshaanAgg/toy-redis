package cmd

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/app/types"
)

func Type(server *types.ServerState, args ...string) []byte {
	if len(args) != 1 {
		return respHandler.Err.Encode(
			fmt.Sprintf("ERR wrong number of arguments for 'TYPE' command: expected 1, got %d", len(args)),
		)
	}

	key := args[0]
	_, okString := server.DB[key]
	_, okStream := server.Streams[key]

	var keyType string
	if !okString && !okStream {
		keyType = "none"
	} else if okString && !okStream {
		keyType = "string"
	} else if okStream && !okString {
		keyType = "stream"
	} else {
		panic("key is both a string and a stream!")
	}

	messageBytes, err := respHandler.Str.Encode(keyType)
	if err != nil {
		fmt.Println("Error encoding response: ", err)
		return nil
	}
	return messageBytes
}
