package cmd

import (
	"fmt"
	"net"
)

func ReplConf(conn net.Conn) {
	bytes, err := respHandler.Str.Encode("OK")
	if err != nil {
		fmt.Println("Failed to encode response", err)
		return
	}
	conn.Write(bytes)
}
