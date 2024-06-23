package cmd

import (
	"fmt"
	"net"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/types"
)

func Incr(conn net.Conn, db *map[string]types.DBItem, key string) {
	val, ok := (*db)[key]
	if !ok {
		(*db)[key] = types.DBItem{
			Value:  "1",
			Expiry: -1,
		}

		_, err := conn.Write(respHandler.Int.Encode(1))
		if err != nil {
			fmt.Printf("Error encoding response: %s\n", err)
		}
		return
	}

	i, err := strconv.Atoi(val.Value)
	if err != nil {
		_, err = conn.Write(respHandler.Err.Encode(("ERR value is not an integer or out of range")))
		if err != nil {
			fmt.Printf("Error encoding response: %s\n", err)
		}
	}

	(*db)[key] = types.DBItem{
		Value:  strconv.Itoa(i + 1),
		Expiry: val.Expiry,
	}
	_, err = conn.Write(respHandler.Int.Encode(i + 1))
	if err != nil {
		fmt.Printf("Error encoding response: %s\n", err)
	}
}
