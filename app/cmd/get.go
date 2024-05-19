package cmd

import (
	"net"
	"time"
)

func Get(con net.Conn, db *map[string]DBItem, key string) {
	value, ok := (*db)[key]
	if !ok {
		res := respHandler.Nil.Encode()
		con.Write(res)
	}

	if value.expiry == -1 || time.Now().UnixMilli() < value.expiry {
		res := respHandler.BulkStr.Encode(value.value)
		con.Write(res)
		return
	}

	delete(*db, key)
	con.Write(respHandler.Nil.Encode())
}
