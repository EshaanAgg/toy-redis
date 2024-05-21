package cmd

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/types"
)

func Wait(conn net.Conn, status *types.ServerState, args ...string) {
	if len(args) != 2 {
		fmt.Printf("Invalid number of arguments for WAIT command: expected 2, got %d\n", len(args))
		return
	}

	numReplicas, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf("Error converting number of replicas to integer: %v\n", err)
		return
	}

	timeOut, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Printf("Error converting timeout to integer: %v\n", err)
		return
	}

	for _, replica := range status.Replicas {
		go func(r types.Replica) {
			err := r.GetAcknowlegment()
			if err != nil {
				fmt.Printf("Error getting acknowledgment from replica: %v\n", err)
			}
		}(replica)
	}

	// Wait for the replicas to acknowledge the bytes
	startTime := time.Now()
	for {
		// Breakout of the loop if the timeout is reached
		if time.Since(startTime) > time.Duration(timeOut)*time.Millisecond {
			fmt.Printf("Timeout reached waiting for %d replicas\n", numReplicas)
			break
		}

		ackCount := getCorrectAckCount(status.Replicas, status.BytesSent)
		if ackCount >= numReplicas {
			break
		}

		time.Sleep(10 * time.Millisecond)
	}

	ackCount := getCorrectAckCount(status.Replicas, status.BytesSent)
	bytes := respHandler.Int.Encode(ackCount)
	_, err = conn.Write(bytes)
	if err != nil {
		fmt.Printf("Error writing number of replicas to connection: %v\n", err)
	}
}

func getCorrectAckCount(replicas []types.Replica, bytesSent int) int {
	ackCount := 0
	for _, replica := range replicas {
		if replica.BytesAcknowledged >= bytesSent {
			ackCount++
		}
	}
	return ackCount
}
