package cmd

import (
	"fmt"
	"net"
)

func Info(con net.Conn, role string) {
	replicationInfo := fmt.Sprintf("role:%s", role)

	res := respHandler.BulkStr.Encode(replicationInfo)
	con.Write(res)
}
