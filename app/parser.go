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

	responseBytes := getCommandResponse(conn, state, arr, buf, isMasterCommand)
	if responseBytes == nil {
		fmt.Println("No response to send to the client")
	} else {
		_, err := conn.Write(responseBytes)
		if err != nil {
			fmt.Printf("Error writing response to client: %v\n", err)
		}
		fmt.Printf("Sent %d bytes to client: %q\n", len(responseBytes), responseBytes)
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

// Executes the command and returns the response in the RESP format
// Returns nil if no response is to be sent to the client
// PSYNC, REPLCONF, WAIT are special commands that are handled by the cmd package directly, and thus their response is nil
func getCommandResponse(conn net.Conn, state *types.ServerState, arr []string, buf []byte, isMasterCommand bool) []byte {
	fmt.Println("Command received: ", arr)

	var res []byte = nil

	switch strings.ToUpper(arr[0]) {
	case "PING":
		res = cmd.Ping(isMasterCommand)

	case "ECHO":
		res = cmd.Echo(arr[1])

	case "SET":
		toReply := !isMasterCommand
		res = cmd.Set(state, toReply, arr[1:]...)
		if state.Role == "master" {
			state.BytesSent += len(buf)
			streamToReplicas(state.Replicas, buf)
		}

	case "GET":
		res = cmd.Get(&state.DB, &state.DBMutex, arr[1])

	case "INCR":
		res = cmd.Incr(&state.DB, arr[1])

	case "INFO":
		res = cmd.Info(state)

	case "REPLCONF":
		cmd.ReplConf(conn, arr[1:], state)

	case "PSYNC":
		cmd.Psync(conn, state.MasterReplID, state.MasterReplOffset)

	case "WAIT":
		cmd.Wait(conn, state, arr[1:]...)

	case "CONFIG":
		res = cmd.Config(state, arr[1:]...)

	case "KEYS":
		res = cmd.Keys(state, arr[1:]...)

	case "TYPE":
		res = cmd.Type(state, arr[1:]...)

	case "XADD":
		res = cmd.Xadd(state, arr[1:]...)

	case "XRANGE":
		res = cmd.Xrange(state, arr[1:]...)

	case "XREAD":
		res = cmd.Xread(state, arr[1:]...)

	default:
		fmt.Printf("Unknown command: %s\n", arr[0])
	}

	return res
}
