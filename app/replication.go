package main

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/types"
)

func sendAndAssertReply(conn net.Conn, messageArr []string, expectedMsg string, respHandler resp.RESPHandler) error {
	bytes := respHandler.Array.Encode(messageArr)
	conn.Write(bytes)

	resp := make([]byte, 1024)
	n, _ := conn.Read(resp)
	msg, err := respHandler.Str.Decode(resp[:n])
	if err != nil {
		return fmt.Errorf("failed to decode response: %s", err)
	}
	if msg != expectedMsg {
		return fmt.Errorf("expected +OK, got %s", string(resp[:n]))
	}

	return nil
}

func sendAndGetReply(conn net.Conn, messageArr []string, respHandler resp.RESPHandler) ([]string, error) {
	bytes := respHandler.Array.Encode(messageArr)
	conn.Write(bytes)

	resp := make([]byte, 1024)
	n, _ := conn.Read(resp)
	msg, _, err := respHandler.Array.Decode(resp[:n])
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %s", err)
	}

	return msg, nil
}

func handshakeWithMaster(server types.ServerState) {
	masterConn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", server.MasterHost, server.MasterPort))
	if err != nil {
		fmt.Println("Failed to connect to master", err)
		return
	}
	defer masterConn.Close()

	respHandler := resp.RESPHandler{}

	// PING
	err = sendAndAssertReply(
		masterConn,
		[]string{"PING"},
		"PONG",
		respHandler,
	)
	if err != nil {
		fmt.Println("Failed to send PING to master", err)
		return
	}

	// REPLCONF listening-port <port>
	err = sendAndAssertReply(
		masterConn,
		[]string{"REPLCONF", "listening-port", fmt.Sprintf("%d", server.Port)},
		"OK",
		respHandler,
	)
	if err != nil {
		fmt.Println("Failed to send REPLCONF listening-port to master", err)
		return
	}

	// REPLCONF capa psync2
	err = sendAndAssertReply(
		masterConn,
		[]string{"REPLCONF", "capa", "psync2"},
		"OK",
		respHandler,
	)
	if err != nil {
		fmt.Println("Failed to send REPLCONF capa psync2 to master", err)
		return
	}

	// PSYNC <replicationid> <offset>
	message, err := sendAndGetReply(
		masterConn,
		[]string{"PSYNC", server.MasterReplID, fmt.Sprintf("%d", server.MasterReplOffset)},
		respHandler,
	)
	if err != nil {
		fmt.Println("Failed to send PSYNC to master", err)
		return
	}
	fmt.Println("Master response to PSYNC:", message)
}
