package main

import (
	"net"
)

// Asserts that the command is a valid Redis command and then calls the appropriate handler
// The bytes should be encoded in the proper RESP format
func handleCommand(buf []byte, conn net.Conn) {

}
