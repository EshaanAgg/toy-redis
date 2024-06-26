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
// isMasterCommand is used to determine if this is a streamed command from the master
func handleCommand(buf []byte, conn net.Conn, state *types.ServerState, isMasterCommand bool) {
	restHandler := resp.RESPHandler{}

	arr, next, err := restHandler.DecodeCommand(buf)
	buf = buf[:len(buf)-len(next)]
	if err != nil {
		fmt.Printf("Error decoding command: %v\n", err)
		return
	}

	fmt.Println("Command received: ", arr)

	switch strings.ToUpper(arr[0]) {

	case "PING":
		cmd.Ping(conn, isMasterCommand)

	case "ECHO":
		cmd.Echo(conn, arr[1])

	case "SET":
		toReply := !isMasterCommand
		cmd.Set(conn, state, toReply, arr[1:]...)
		if state.Role == "master" {
			state.BytesSent += len(buf)
			streamToReplicas(state.Replicas, buf)
		}

	case "GET":
		cmd.Get(conn, &state.DB, &state.DBMutex, arr[1])

	case "INFO":
		cmd.Info(conn, state)

	case "REPLCONF":
		cmd.ReplConf(conn, arr[1:], state)

	case "PSYNC":
		cmd.Psync(conn, state.MasterReplID, state.MasterReplOffset)

	case "WAIT":
		cmd.Wait(conn, state, arr[1:]...)

	case "CONFIG":
		cmd.Config(conn, state, arr[1:]...)

	case "KEYS":
		cmd.Keys(conn, state, arr[1:]...)

	case "TYPE":
		cmd.Type(conn, state, arr[1:]...)

	case "XADD":
		cmd.Xadd(conn, state, arr[1:]...)

	case "XRANGE":
		cmd.Xrange(conn, state, arr[1:]...)

	case "XREAD":
		cmd.Xread(conn, state, arr[1:]...)

	default:
		fmt.Printf("Unknown command: %s\n", arr[0])
	}

	// If this was a command from master, update the acknowledgment offset
	if isMasterCommand {
		state.AckOffset += len(buf)
	}

	// If there are more commands in the buffer, handle them
	if len(next) > 0 {
		handleCommand(next, conn, state, isMasterCommand)
	}
}
