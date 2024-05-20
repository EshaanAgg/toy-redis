package cmd

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/types"
)

func Wait(conn net.Conn, status *types.ServerState, args ...string) {
	if len(args) != 2 {
		fmt.Printf("Invalid number of arguments for WAIT command: expected 2, got %d\n", len(args))
		return
	}

	// numReplicas, err := strconv.Atoi(args[0])
	// if err != nil {
	// 	fmt.Printf("Error converting number of replicas to integer: %v\n", err)
	// 	return
	// }
	bytes := respHandler.Int.Encode(len(status.ReplicaConn))

	_, err := conn.Write(bytes)
	if err != nil {
		fmt.Printf("Error writing number of replicas to connection: %v\n", err)
	}
}
