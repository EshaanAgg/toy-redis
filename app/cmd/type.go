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
	_, ok := server.DB[key]

	var keyType string
	if !ok {
		keyType = "none"
	} else {
		keyType = "string"
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
