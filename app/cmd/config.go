package cmd

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/app/types"
)

func Config(status *types.ServerState, args ...string) []byte {
	if len(args) != 2 {
		return respHandler.Err.Encode(
			fmt.Sprintf("ERR wrong number of arguments for 'CONFIG' command: expected 2, got %d", len(args)),
		)
	}

	if args[0] == "GET" {
		if args[1] == "dir" {
			return respHandler.Array.Encode([]string{"dir", status.DBDir})
		}

		if args[1] == "dbfilename" {
			return respHandler.Array.Encode([]string{"dbfilename", status.DBFilename})
		}

		return respHandler.Err.Encode(
			fmt.Sprintf("ERR unknown CONFIG GET command: %s\n", args),
		)
	}

	return respHandler.Err.Encode(
		fmt.Sprintf("ERR unknown CONFIG command: %s\n", args),
	)
}
