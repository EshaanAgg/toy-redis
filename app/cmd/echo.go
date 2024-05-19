package cmd

import (
	"net"
)

func Echo(con net.Conn, message string) {
	res := respHandler.BulkStr.Encode(message)
	con.Write(res)
}
