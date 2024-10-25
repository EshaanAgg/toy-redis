package cmd

import (
	"fmt"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/types"
)

func Incr(db *map[string]types.DBItem, key string) []byte {
	val, ok := (*db)[key]
	if !ok {
		(*db)[key] = types.DBItem{
			Value:  "1",
			Expiry: -1,
		}

		return respHandler.Int.Encode(1)
	}

	i, err := strconv.Atoi(val.Value)
	if err != nil {
		return respHandler.Err.Encode(("ERR value is not an integer or out of range"))
	}

	(*db)[key] = types.DBItem{
		Value:  fmt.Sprint(i + 1),
		Expiry: val.Expiry,
	}
	return respHandler.Int.Encode(i + 1)
}
