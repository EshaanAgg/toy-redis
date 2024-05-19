package main

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/cmd"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

// Asserts that the command is a valid Redis command and then calls the appropriate handler
// The bytes should be encoded in the proper RESP format
func handleCommand(buf []byte, conn net.Conn) {
	restHandler := resp.RESPHandler{}

	arr, err := restHandler.DecodeCommand(buf)
	if err != nil {
		fmt.Printf("Error decoding command: %v\n", err)
		return
	}

	switch arr[0] {
	case "PING":
		cmd.Ping(conn)
	case "ECHO":
		cmd.Echo(conn, arr[1])
	default:
		fmt.Printf("Unknown command: %s\n", arr[0])
	}
}
