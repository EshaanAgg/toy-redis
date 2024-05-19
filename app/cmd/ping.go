package cmd

import (
	"fmt"
	"net"
)

func Ping(con net.Conn) {
	res, err := respHandler.Str.Encode("PONG")

	if err != nil {
		fmt.Printf("Error encoding response: %s\n", err)
		return
	}

	con.Write(res)
}
