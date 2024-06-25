package cmd

import (
	"fmt"
	"net"

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
	server.Transactions[conn] = types.TransactionData{
		Started: true,
		Queue:   [][]byte{},
	}
	server.TransactionMutex.Unlock()
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

func executeTransaction(t *types.TransactionData, server *types.ServerState) []byte {
	output := make([][]byte, 0)

	for _, command := range t.Queue {
		arr, next, err := respHandler.DecodeCommand(command)
		if err != nil {
			fmt.Printf("Error decoding queued command '%q': %v\n", command, err)
			return nil
		}
		if len(next) > 0 {
			fmt.Printf("Multiple commands have been queued together, which is not expected, bytes: %q\n", command)
			return nil
		}

		fmt.Println("Executing queued command: ", arr[0])

		responseBytes := GetCommandResponse(nil, server, arr, command, false)
		if responseBytes == nil {
			fmt.Printf("No response to send to the client for command: %s\n", arr[0])
			return nil
		}
		fmt.Printf("Queued command response: %q\n", responseBytes)
		output = append(output, responseBytes)
	}

	return respHandler.Array.EncodeFromElementBytes(output)
}
