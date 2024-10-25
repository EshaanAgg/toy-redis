package cmd

import (
	"fmt"
	"net"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/types"
)

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
		sendOk(conn)

	case "EXEC":
		if !IsTransactionStarted(conn, server) {
			sendResponse(
				conn,
				respHandler.Err.Encode("ERR EXEC without MULTI"),
			)
			return
		}

		transaction, ok := server.Transactions[conn]
		if !ok {
			panic("Transaction not found but is started")
		}
		res := executeTransaction(&transaction, server)
		if res == nil {
			fmt.Println("There was some error executing the transaction and no response is present to send to the client")
			return
		}
		endTransaction(conn, server)
		sendResponse(conn, res)

	case "DISCARD":
		if !IsTransactionStarted(conn, server) {
			sendResponse(
				conn,
				respHandler.Err.Encode("ERR DISCARD without MULTI"),
			)
			return
		}

		endTransaction(conn, server)
		sendOk(conn)

	default:
		panic(fmt.Sprintf("Unknown transaction command: %s", command))
	}
}

func sendResponse(conn net.Conn, response []byte) {
	_, err := conn.Write(response)
	if err != nil {
		fmt.Printf("Error writing response to client: %v\n", err)
	}
	fmt.Printf("Sent %d bytes to client: %q\n", len(response), response)
}
