package file

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/redis-starter-go/app/types"
)

func InitialiseDB(state *types.ServerState, dbFile string, dbPath string) {
	filePath := fmt.Sprintf("%s/%s", dbPath, dbFile)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// File does not exist, create a new file
		fmt.Printf("No file exists at the provided path: %s\nInitialized with empty database\n", filePath)
		return
	}

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading database file from disk: %s\n", err)
		return
	}

	err = parseFile(fileContent, state)
	if err != nil {
		fmt.Printf("Error parsing database file: %s\n", err)
	}
}
