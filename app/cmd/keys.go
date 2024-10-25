package cmd

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/app/types"
)

func Keys(server *types.ServerState, args ...string) []byte {
	if len(args) != 1 || args[0] != "*" {
		return respHandler.Err.Encode(
			fmt.Sprintf("ERR invalid number of arguments for 'KEYS' command: %v\n", args),
		)
	}

	keys := make([]string, 0)
	for key := range server.DB {
		keys = append(keys, key)
	}

	return respHandler.Array.Encode(keys)
}
