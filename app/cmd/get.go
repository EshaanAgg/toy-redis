package cmd

import (
	"sync"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/types"
)

func Get(db *map[string]types.DBItem, mutex *sync.Mutex, key string) []byte {
	mutex.Lock()
	defer mutex.Unlock()

	value, ok := (*db)[key]

	if !ok {
		return respHandler.Nil.Encode()
	}

	if value.Expiry == -1 || time.Now().UnixMilli() < value.Expiry {
		return respHandler.BulkStr.Encode(value.Value)
	}

	delete(*db, key)
	return respHandler.Nil.Encode()
}
