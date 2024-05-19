package cmd

import (
	"net"
)

func Get(con net.Conn, db *map[string]string, key string) {
	value, ok := (*db)[key]

	if !ok {
		res := respHandler.Nil.Encode()
		con.Write(res)
	}

	res := respHandler.BulkStr.Encode(value)
	con.Write(res)
}
