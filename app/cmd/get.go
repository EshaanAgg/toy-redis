package cmd

import (
	"net"
	"sync"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/types"
)

func Get(con net.Conn, db *map[string]types.DBItem, mutex *sync.Mutex, key string) {
	mutex.Lock()
	defer mutex.Unlock()

	value, ok := (*db)[key]

	if !ok {
		res := respHandler.Nil.Encode()
		con.Write(res)
		return
	}

	if value.Expiry == -1 || time.Now().UnixMilli() < value.Expiry {
		res := respHandler.BulkStr.Encode(value.Value)
		con.Write(res)
		return
	}

	delete(*db, key)
	con.Write(respHandler.Nil.Encode())
}
