package types

import (
	"net"
	"sync"
)

type DBItem struct {
	Value  string
	Expiry int64
}

type ServerState struct {
	DB      map[string]DBItem
	DBMutex sync.Mutex
	Port    int

	Role             string      // master | slave
	MasterReplID     string      // Replication ID of the master (own replication ID if master)
	MasterReplOffset int         // Offset of the master (0 if master)
	MasterHost       string      // Host of the master (empty if master)
	MasterPort       string      // Port of the master (empty if master)
	ReplicaConn      []*net.Conn // Connections to replicas (empty if slave)
	AckOffset        int         // Offset of the last acknowledged replication message
}
