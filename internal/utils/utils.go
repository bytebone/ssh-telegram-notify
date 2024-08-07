package utils

import (
	"log"
	"os"
)

// GetHomeDir returns the home directory of the current user.
// When an error is encountered, the code is fataled.
func GetHomeDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	return homeDir
}
