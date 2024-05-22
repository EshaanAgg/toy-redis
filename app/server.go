package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/file"
	"github.com/codecrafters-io/redis-starter-go/app/types"
)

func NewServerState(args *Args) *types.ServerState {
	state := types.ServerState{
		DB:   map[string]types.DBItem{},
		Port: args.port,

		Role:             "master",
		MasterReplID:     "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb",
		MasterReplOffset: 0,

		DBDir:      args.dir,
		DBFilename: args.dbfilename,
	}

	if args.replicaof != "" {
		state.Role = "slave"
		state.MasterHost = strings.Split(args.replicaof, " ")[0]
		state.MasterPort = strings.Split(args.replicaof, " ")[1]
		state.MasterReplID = "?"
		state.MasterReplOffset = -1
		handshakeWithMaster(&state)
	}

	if args.dir != "" && args.dbfilename != "" {
		file.InitialiseDB(&state, args.dbfilename, args.dir)
	}

	// Initialise the streams map
	state.Streams = map[string][]types.StreamEntry{}

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
		fmt.Printf("Accepted connection from %s\n", conn.LocalAddr().String())

		go handleConnection(conn, serverState, false)
	}
}

func handleConnection(conn net.Conn, state *types.ServerState, isMasterConnection bool) {
	defer conn.Close()

	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Printf("Connection closed to %s\n", conn.LocalAddr().String())
				return
			}
			fmt.Println("Error reading:", err)
			return
		}

		fmt.Printf("Received %d bytes: %q\n", n, buf[:n])
		handleCommand(buf[:n], conn, state, isMasterConnection)
	}
}
