package cmd

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/types"
)

func Keys(conn net.Conn, server *types.ServerState, args ...string) {
	if len(args) != 1 || args[0] != "*" {
		fmt.Printf("Invalid number of arguments for 'KEYS' command: %v\n", args)
		return
	}

	keys := make([]string, 0)
	for key := range server.DB {
		keys = append(keys, key)
	}

	messageBytes := respHandler.Array.Encode(keys)
	_, err := conn.Write(messageBytes)
	if err != nil {
		fmt.Println("Error writing response to connection: ", err)
	}
}
