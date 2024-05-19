package main

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/types"
)

func handshakeWithMaster(server types.ServerState) {
	masterConn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", server.MasterHost, server.MasterPort))
	if err != nil {
		fmt.Println("Failed to connect to master", err)
		return
	}
	defer masterConn.Close()

	respHandler := resp.RESPHandler{}

	// PING
	pingArray := respHandler.Array.Encode([]string{"PING"})
	masterConn.Write(pingArray)

	// REPLCONF listening-port <port>
	replconfArray := respHandler.Array.Encode([]string{"REPLCONF", "listening-port", server.MasterPort})
	masterConn.Write(replconfArray)

	// REPLCONF capa psync2
	replconfArray = respHandler.Array.Encode([]string{"REPLCONF", "capa", "psyc2"})
	masterConn.Write(replconfArray)
}
