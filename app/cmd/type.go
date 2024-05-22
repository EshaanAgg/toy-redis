package cmd

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/types"
)

func Type(conn net.Conn, server *types.ServerState, args ...string) {
	if len(args) != 1 {
		fmt.Printf("Invalid number of arguments for 'TYPE' command: %v\n", args)
		return
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
		panic("Key is both a string and a stream!")
	}

	messageBytes, err := respHandler.Str.Encode(keyType)
	if err != nil {
		fmt.Println("Error encoding response: ", err)
		return
	}
	_, err = conn.Write(messageBytes)
	if err != nil {
		fmt.Println("Error writing response to connection: ", err)
	}
}
