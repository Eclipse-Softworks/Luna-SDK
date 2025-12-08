package main

import (
	"os"

	"github.com/eclipse-softworks/luna-sdk/cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
