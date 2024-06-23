package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func main() {
	port := flag.Int("port", 3000, "The port on which the Redis server listens to")
	flag.Parse()

	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		fmt.Printf("Couldn't establish connection to server: %v", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected successfully to the Redis server.")
	respHandler := resp.RESPHandler{}

	// Start the REPL and handle communication with the server
	for {
		fmt.Printf("> ")

		in := bufio.NewReader(os.Stdin)
		cmd, err := in.ReadString('\n')
		if err != nil {
			fmt.Printf("There was an error in reading the command from the STDIN: %v", err)
			continue
		}

		cmdParts := strings.Split(cmd[:len(cmd)-1], " ")
		_, err = conn.Write(respHandler.Array.Encode(cmdParts))
		if err != nil {
			fmt.Printf("There was an error in sending the message: %v\n", err)
			continue
		}

		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Printf("There was an error in recieving the message: %v\n", err)
			continue
		}

		res, err := respHandler.DecodeResponse(buffer[:n])
		if err != nil {
			fmt.Printf("There was an error in decoding the response: %v\n", err)
			continue
		}

		fmt.Printf(">> %s\n", res)
	}
}
