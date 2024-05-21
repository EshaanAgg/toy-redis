package types

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type Replica struct {
	Conn              net.Conn
	BytesAcknowledged int
}

func (r *Replica) GetAcknowlegment() error {
	respHandler := resp.RESPHandler{}

	// Send the GETACK command to the replica
	messageBytes := respHandler.Array.Encode([]string{"REPLCONF", "GETACK", "*"})
	_, err := r.Conn.Write(messageBytes)
	if err != nil {
		return fmt.Errorf("failed to write to replica connection: %v", err)
	}

	return nil
}
