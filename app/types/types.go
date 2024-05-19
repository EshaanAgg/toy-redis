package types

import "net"

type DBItem struct {
	Value  string
	Expiry int64
}

type ServerState struct {
	DB   map[string]DBItem
	Port int

	Role             string      // master | slave
	MasterReplID     string      // Replication ID of the master (own replication ID if master)
	MasterReplOffset int         // Offset of the master (0 if master)
	MasterHost       string      // Host of the master (empty if master)
	MasterPort       string      // Port of the master (empty if master)
	ReplicaConn      []*net.Conn // Connections to replicas (empty if slave)
}
