package main

import (
	"os"

	"github.com/dreamlibrarian/solaredge-monitoring/cmd"
)

func main() {
	err := cmd.RootCmd.Execute()
	if err != nil {
		os.Exit(1)
		return
	}
}
