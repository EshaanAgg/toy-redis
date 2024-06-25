package main

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/cmd"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/types"
)

// Asserts that the command is a valid Redis command and then calls the appropriate handler
// The bytes should be encoded in the proper RESP format
// isMasterCommand is used to determine if this is a streamed command from the master
func handleCommand(buf []byte, conn net.Conn, state *types.ServerState, isMasterCommand bool) {
	restHandler := resp.RESPHandler{}

	arr, next, err := restHandler.DecodeCommand(buf)
	buf = buf[:len(buf)-len(next)]
	if err != nil {
		fmt.Printf("Error decoding command: %v\n", err)
		return
	}

	// Check if the command is a transaction command and handle it
	if cmd.IsTransactionCommand(arr[0]) {
		cmd.HandleTransactionCommand(conn, arr[0], state)
	} else {
		inTransaction := cmd.IsTransactionStarted(conn, state)

		if inTransaction {
			// If in a transaction, queue the command
			// and send a QUEUED response to the client

			state.TransactionMutex.Lock()
			state.Transactions[conn] = types.TransactionData{
				Started: true,
				Queue:   append(state.Transactions[conn].Queue, buf),
			}
			state.TransactionMutex.Unlock()

			res, err := resp.RESPHandler{}.Str.Encode("QUEUED")
			if err != nil {
				fmt.Printf("Error encoding response: %v\n", err)
			}
			_, err = conn.Write(res)
			if err != nil {
				fmt.Printf("Error writing response to client: %v\n", err)
			}

		} else {
			// If not in a transaction, execute the command
			// and send the response to the client

			responseBytes := cmd.GetCommandResponse(conn, state, arr, buf, isMasterCommand)
			if responseBytes == nil {
				fmt.Println("No response to send to the client")
			} else {
				// Send the response to the client if not in a transaction
				_, err := conn.Write(responseBytes)
				if err != nil {
					fmt.Printf("Error writing response to client: %v\n", err)
				}
				fmt.Printf("Sent %d bytes to client: %q\n", len(responseBytes), responseBytes)
			}

		}
	}

	// If this was a command from master, update the acknowledgment offset
	if isMasterCommand {
		state.AckOffset += len(buf)
	}

	// If there are more commands in the buffer, handle them
	if len(next) > 0 {
		handleCommand(next, conn, state, isMasterCommand)
	}
}
