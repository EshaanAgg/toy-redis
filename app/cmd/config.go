package cmd

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/types"
)

func Config(conn net.Conn, status *types.ServerState, args ...string) {
	if len(args) != 2 {
		fmt.Printf("Invalid number of arguments for CONFIG command: expected 2, got %d\n", len(args))
	}

	if args[0] == "GET" {
		if args[1] == "dir" {
			mess := respHandler.Array.Encode([]string{"dir", status.DBDir})
			_, err := conn.Write(mess)
			if err != nil {
				fmt.Printf("Error writing the database directory to connection: %v\n", err)
			}
			return
		}

		if args[1] == "dbfilename" {
			mess := respHandler.Array.Encode([]string{"dbfilename", status.DBFilename})
			_, err := conn.Write(mess)
			if err != nil {
				fmt.Printf("Error writing the database filename to connection: %v\n", err)
			}
			return
		}

		fmt.Printf("Unknown CONFIG GET command: %s\n", args)
		return
	}

	fmt.Printf("Unknown CONFIG command: %s\n", args)
}
