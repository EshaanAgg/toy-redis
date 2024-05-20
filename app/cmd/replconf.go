package cmd

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/types"
)

func ReplConf(conn net.Conn, args []string, serverState *types.ServerState) {
	// If the command is REPLCONF listening-port, add the connection to the list of replica connections
	if len(args) >= 2 && args[0] == "listening-port" {
		serverState.ReplicaConn = append(serverState.ReplicaConn, &conn)
		sendOk(conn)
		return
	}

	// If the command is REPLCONF capa, send OK
	if len(args) >= 2 && args[0] == "capa" {
		sendOk(conn)
		return
	}

	// If the command is REPLCONF GETACK *, send an ACK back to the replica
	if len(args) >= 2 && args[0] == "GETACK" && args[1] == "*" {
		getAck(conn, serverState.AckOffset)
		return
	}

	fmt.Printf("Unknown REPLCONF command: %s\n", args)
}

func getAck(conn net.Conn, bytesOffset int) {
	bytes := respHandler.Array.Encode([]string{"REPLCONF", "ACK", fmt.Sprintf("%d", bytesOffset)})
	_, err := conn.Write(bytes)
	if err != nil {
		fmt.Println("Failed to write response ACK response to master: ", err)
	}
}

func sendOk(conn net.Conn) {
	bytes, err := respHandler.Str.Encode("OK")
	if err != nil {
		fmt.Println("Failed to encode OK response: ", err)
		return
	}
	_, err = conn.Write(bytes)
	if err != nil {
		fmt.Println("Failed to write OK response to master: ", err)
	}
}
