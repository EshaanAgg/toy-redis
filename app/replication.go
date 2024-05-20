package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"

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

func sendAndGetRBDFile(conn net.Conn, messageArr []string, respHandler resp.RESPHandler, state *types.ServerState) (string, error) {
	bytes := respHandler.Array.Encode(messageArr)
	conn.Write(bytes)

	// Get the initial PSYNC response
	resp := make([]byte, 1024)
	n, err := conn.Read(resp)
	if err != nil {
		return "", fmt.Errorf("failed to recieve message from master: %s", err)
	}
	psyncResp, err := respHandler.Str.Decode(resp[:n])
	if err != nil {
		return "", fmt.Errorf("failed to decode response: %s", err)
	}

	// Parse the PSYNC response
	responseParts := strings.Split(psyncResp, " ")
	if len(responseParts) != 3 {
		return "", fmt.Errorf("expected 3 parts in PSYNC response, got %d", len(responseParts))
	}
	if responseParts[0] != "FULLRESYNC" {
		return "", fmt.Errorf("expected FULLRESYNC in PSYNC response, got %s", responseParts[0])
	}
	state.MasterReplID = responseParts[1]
	portAsInt, err := strconv.Atoi(responseParts[2])
	if err != nil {
		return "", fmt.Errorf("failed to convert port to int: %s", err)
	}
	state.MasterReplOffset = portAsInt

	// Get the RDB file
	rdbBytes := make([]byte, 1024)
	n, err = conn.Read(rdbBytes)
	if err != nil {
		return "", fmt.Errorf("failed to recieve message from master: %s", err)
	}
	fileContent := string(rdbBytes[:n])

	return fileContent, nil
}

func handshakeWithMaster(server *types.ServerState) {
	masterConn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", server.MasterHost, server.MasterPort))
	if err != nil {
		fmt.Println("Failed to connect to master: ", err)
		return
	}

	respHandler := resp.RESPHandler{}

	// PING
	err = sendAndAssertReply(
		masterConn,
		[]string{"PING"},
		"PONG",
		respHandler,
	)
	if err != nil {
		fmt.Println("Failed to send PING to master: ", err)
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
		fmt.Println("Failed to send REPLCONF listening-port to master: ", err)
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
		fmt.Println("Failed to send REPLCONF capa psync2 to master: ", err)
		return
	}

	// PSYNC <replicationid> <offset>
	rdbFile, err := sendAndGetRBDFile(
		masterConn,
		[]string{"PSYNC", "?", fmt.Sprintf("%d", -1)},
		respHandler,
		server,
	)
	if err != nil {
		fmt.Println("Failed to send PSYNC to master: ", err)
		return
	}
	fmt.Printf("RDB File content: %q\n", rdbFile)

	// Since the handshake was successful, we can now set handle the master connection in a separate goroutine
	go handleConnection(masterConn, server, true)
}

func streamToReplicas(replicaConn []*net.Conn, buff []byte) {
	fmt.Printf("Streaming recieved command to %d replicas\n", len(replicaConn))
	for _, conn := range replicaConn {
		_, err := (*conn).Write(buff)
		if err != nil {
			fmt.Printf("Failed to stream to replica %s: %s", (*conn).RemoteAddr().String(), err.Error())
		}
	}
}
