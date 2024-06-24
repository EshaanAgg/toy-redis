package cmd

import (
	"fmt"
	"net"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/types"
)

func IsTransactionStarted(conn net.Conn, server *types.ServerState) bool {
	server.TransactionMutex.Lock()
	defer server.TransactionMutex.Unlock()

	transaction, ok := server.Transactions[conn]
	return ok && transaction.Started
}

func startTransaction(conn net.Conn, server *types.ServerState) {
	server.TransactionMutex.Lock()
	defer server.TransactionMutex.Unlock()

	server.Transactions[conn] = types.TransactionData{
		Started: true,
		Queue:   [][]byte{},
	}
}

func endTransaction(conn net.Conn, server *types.ServerState) {
	server.TransactionMutex.Lock()
	server.Transactions[conn] = types.TransactionData{
		Started: false,
		Queue:   [][]byte{},
	}
	server.TransactionMutex.Unlock()
}

func IsTransactionCommand(command string) bool {
	return command == "MULTI" || command == "EXEC" || command == "DISCARD"
}

func HandleTransactionCommand(conn net.Conn, command string, server *types.ServerState) {
	switch strings.ToUpper(command) {
	case "MULTI":
		if IsTransactionStarted(conn, server) {
			sendResponse(
				conn,
				respHandler.Err.Encode("ERR MULTI calls can not be nested"),
			)
			endTransaction(conn, server)
			return
		}

		startTransaction(conn, server)
		ok, err := respHandler.Str.Encode("OK")
		if err != nil {
			fmt.Printf("Error encoding response: %v\n", err)
		}
		sendResponse(conn, ok)

	case "EXEC":
	case "DISCARD":
	default:
		panic(fmt.Sprintf("Unknown transaction command: %s", command))
	}
}

func sendResponse(conn net.Conn, response []byte) {
	_, err := conn.Write(response)
	if err != nil {
		fmt.Printf("Error writing response to client: %v\n", err)
	}
}
