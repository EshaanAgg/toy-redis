package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/app/cmd"
)

type Args struct {
	port int
}

func GetArgs() Args {
	port := flag.Int("port", 6379, "The port on which the Redis server listens")
	flag.Parse()

	return Args{
		port: *port,
	}
}

func main() {
	args := GetArgs()

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

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	buf := make([]byte, 1024)
	db := map[string]cmd.DBItem{}

	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err)
			return
		}

		fmt.Printf("Received %d bytes: %s\n", n, buf[:n])
		handleCommand((buf[:n]), conn, db)
	}
}
