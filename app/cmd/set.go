package cmd

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/types"
)

func Set(con net.Conn, db *map[string]types.DBItem, shouldReply bool, arr ...string) {
	if len(arr) < 2 {
		fmt.Println("Error: SET requires at least 2 arguments, which are the KEY and the VALUE")
		return
	}

	expiry := int64(-1)
	if len(arr) > 2 {
		if arr[2] != "px" {
			fmt.Println("Error: SET only supports px as a third argument")
			return
		}
		n, err := strconv.ParseInt(arr[3], 10, 64)
		if err != nil {
			fmt.Println("Error: EX argument must be an integer")
			return
		}
		expiry = time.Now().UnixMilli() + n
	}

	key := arr[0]
	value := arr[1]

	(*db)[key] = types.DBItem{Value: value, Expiry: expiry}

	if shouldReply {
		res, err := respHandler.Str.Encode("OK")
		if err != nil {
			fmt.Printf("Error encoding response: %s\n", err)
			return
		}
		con.Write(res)
	}
}
