package main

import (
	"fmt"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/app/types"
)

func NewServerState(args *Args) *types.ServerState {
	state := types.ServerState{
		DB:               map[string]types.DBItem{},
		Role:             "master",
		MasterReplID:     "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb",
		MasterReplOffset: 0,
	}

	if args.replicaof != "" {
		state.Role = "slave"
	}
	return &state
}

func main() {
	args := GetArgs()
	serverState := NewServerState(&args)

	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", args.port))
	if err != nil {
		fmt.Println("Failed to bind to port 6379", err)
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err)
			os.Exit(1)
		}
		fmt.Printf("Accepted connection from %s\n", conn.RemoteAddr().String())

		go handleConnection(conn, serverState)
	}
}

func handleConnection(conn net.Conn, state *types.ServerState) {
	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err)
			return
		}

		fmt.Printf("Received %d bytes: %s\n", n, buf[:n])
		handleCommand(buf[:n], conn, state)
	}
}
