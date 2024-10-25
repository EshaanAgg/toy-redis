package cmd

import (
	"fmt"
)

func Ping(isMasterCommand bool) []byte {
	if isMasterCommand {
		// If the PING command is from the master, we should not reply
		// The master is only sending PING to check if the replica is alive
		return nil
	}

	res, err := respHandler.Str.Encode("PONG")
	if err != nil {
		fmt.Printf("Error encoding response: %s\n", err)
		return nil
	}
	return res
}
