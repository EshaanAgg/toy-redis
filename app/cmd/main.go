package cmd

import (
	"fmt"
	"net"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/types"
)

var respHandler = resp.RESPHandler{}

// Executes the command and returns the response in the RESP format
// Returns nil if no response is to be sent to the client
// PSYNC, REPLCONF, WAIT are special commands that are handled by the cmd package directly, and thus their response is nil
func GetCommandResponse(conn net.Conn, state *types.ServerState, arr []string, buf []byte, isMasterCommand bool) []byte {
	fmt.Println("Command received: ", arr)

	var res []byte = nil

	switch strings.ToUpper(arr[0]) {
	case "PING":
		res = Ping(isMasterCommand)

	case "ECHO":
		res = Echo(arr[1])

	case "SET":
		toReply := !isMasterCommand
		res = Set(state, toReply, arr[1:]...)
		if state.Role == "master" {
			state.BytesSent += len(buf)
			streamToReplicas(state.Replicas, buf)
		}

	case "GET":
		res = Get(&state.DB, &state.DBMutex, arr[1])

	case "INCR":
		res = Incr(&state.DB, arr[1])

	case "INFO":
		res = Info(state)

	case "REPLCONF":
		ReplConf(conn, arr[1:], state)

	case "PSYNC":
		Psync(conn, state.MasterReplID, state.MasterReplOffset)

	case "WAIT":
		Wait(conn, state, arr[1:]...)

	case "CONFIG":
		res = Config(state, arr[1:]...)

	case "KEYS":
		res = Keys(state, arr[1:]...)

	case "TYPE":
		res = Type(state, arr[1:]...)

	case "XADD":
		res = Xadd(state, arr[1:]...)

	case "XRANGE":
		res = Xrange(state, arr[1:]...)

	case "XREAD":
		res = Xread(state, arr[1:]...)

	default:
		fmt.Printf("Unknown command: %s\n", arr[0])
	}

	return res
}

func streamToReplicas(replicas []types.Replica, buff []byte) {
	fmt.Printf("Streaming recieved command to %d replicas\n", len(replicas))
	for ind, r := range replicas {
		_, err := r.Conn.Write(buff)
		if err != nil {
			fmt.Printf("Failed to stream to replica %d: %s", ind+1, err.Error())
		}
	}
}
