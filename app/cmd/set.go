package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/types"
)

func Set(server *types.ServerState, shouldReply bool, arr ...string) []byte {
	if len(arr) < 2 {
		return respHandler.Err.Encode(
			fmt.Sprintf("ERR wrong number of arguments for '%s' command", arr[0]),
		)
	}

	expiry := int64(-1)
	if len(arr) > 2 {
		if arr[2] != "px" {
			return respHandler.Err.Encode("ERR set only supports px as a third argument")
		}
		n, err := strconv.ParseInt(arr[3], 10, 64)
		if err != nil {
			return respHandler.Err.Encode("ERR EX argument must be an integer")
		}
		expiry = time.Now().UnixMilli() + n
	}

	server.DBMutex.Lock()
	defer server.DBMutex.Unlock()

	key := arr[0]
	value := arr[1]

	server.DB[key] = types.DBItem{Value: value, Expiry: expiry}

	if shouldReply {
		res, err := respHandler.Str.Encode("OK")
		if err != nil {
			fmt.Printf("Error encoding response: %s\n", err)
			return nil
		}
		return res
	}

	return nil
}
