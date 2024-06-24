package cmd

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/app/types"
)

func Info(serverInfo *types.ServerState) []byte {
	replicationInfo := fmt.Sprintf("role:%s", serverInfo.Role)
	replicationInfo += fmt.Sprintf("\nmaster_replid:%s", serverInfo.MasterReplID)
	replicationInfo += fmt.Sprintf("\nmaster_repl_offset:%d", serverInfo.MasterReplOffset)

	return respHandler.BulkStr.Encode(replicationInfo)
}
