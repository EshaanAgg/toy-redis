package cmd

import (
	"net"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/types"
)

func Get(con net.Conn, db *map[string]types.DBItem, key string) {
	value, ok := (*db)[key]
	if !ok {
		res := respHandler.Nil.Encode()
		con.Write(res)
	}

	if value.Expiry == -1 || time.Now().UnixMilli() < value.Expiry {
		res := respHandler.BulkStr.Encode(value.Value)
		con.Write(res)
		return
	}

	delete(*db, key)
	con.Write(respHandler.Nil.Encode())
}
