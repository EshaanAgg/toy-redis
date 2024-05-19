package cmd

import "net"

func Ping(con net.Conn) {
	con.Write([]byte("+PONG\r\n"))
}
