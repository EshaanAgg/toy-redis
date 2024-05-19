package cmd

import "net"

func Info(con net.Conn) {
	res := respHandler.BulkStr.Encode("role:master")
	con.Write(res)
}
