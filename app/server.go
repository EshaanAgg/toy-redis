package main

import (
	"fmt"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/app/cmd"
)

type ServerState struct {
	db   map[string]cmd.DBItem
	role string
}

func NewServerState(args *Args) *ServerState {
	state := ServerState{
		db:   map[string]cmd.DBItem{},
		role: "master",
	}

	if args.replicaof != "" {
		state.role = "slave"
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

func handleConnection(conn net.Conn, state *ServerState) {
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
