package cmd

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/types"
)

func Info(con net.Conn, serverInfo *types.ServerState) {
	replicationInfo := fmt.Sprintf("role:%s", serverInfo.Role)
	replicationInfo += fmt.Sprintf("\nmaster_replid:%s", serverInfo.MasterReplID)
	replicationInfo += fmt.Sprintf("\nmaster_repl_offset:%d", serverInfo.MasterReplOffset)

	res := respHandler.BulkStr.Encode(replicationInfo)
	con.Write(res)
}
