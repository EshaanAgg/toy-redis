package cmd

import (
	"fmt"
	"net"
)

func Psync(conn net.Conn, replID string, offset int) {
	bytes, err := respHandler.Str.Encode(fmt.Sprintf("FULLRESYNC %s %d", replID, offset))
	if err != nil {
		fmt.Println("Failed to encode response", err)
		return
	}
	conn.Write(bytes)
}
