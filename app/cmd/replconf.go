package cmd

import (
	"fmt"
	"net"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/types"
)

func ReplConf(conn net.Conn, args []string, serverState *types.ServerState) {
	bytes, err := respHandler.Str.Encode("OK")
	if err != nil {
		fmt.Println("Failed to encode response", err)
		return
	}
	conn.Write(bytes)

	if len(args) >= 2 && args[0] == "listening-port" {
		port, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Failed to parse port from the replica: ", err)
			return
		}
		serverState.ReplicaPorts = append(serverState.ReplicaPorts, port)
	}
}
