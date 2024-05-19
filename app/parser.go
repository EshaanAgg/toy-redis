package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/cmd"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/types"
)

// Asserts that the command is a valid Redis command and then calls the appropriate handler
// The bytes should be encoded in the proper RESP format
func handleCommand(buf []byte, conn net.Conn, state *types.ServerState) {
	restHandler := resp.RESPHandler{}

	arr, err := restHandler.DecodeCommand(buf)
	if err != nil {
		fmt.Printf("Error decoding command: %v\n", err)
		return
	}

	switch strings.ToUpper(arr[0]) {
	case "PING":
		cmd.Ping(conn)
	case "ECHO":
		cmd.Echo(conn, arr[1])
	case "SET":
		cmd.Set(conn, &state.DB, arr[1:]...)
	case "GET":
		cmd.Get(conn, &state.DB, arr[1])
	case "INFO":
		cmd.Info(conn, state)
	case "REPLCONF":
		cmd.ReplConf(conn)
	case "PSYNC":
		cmd.Psync(conn, state.MasterReplID, state.MasterReplOffset)
	default:
		fmt.Printf("Unknown command: %s\n", arr[0])
	}
}
