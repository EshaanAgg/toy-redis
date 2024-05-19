package cmd

import (
	"fmt"
	"net"
)

func Set(con net.Conn, db *map[string]string, key string, value string) {
	(*db)[key] = value
	res, err := respHandler.Str.Encode("OK")
	if err != nil {
		fmt.Printf("Error encoding response: %s\n", err)
		return
	}
	con.Write(res)
}
