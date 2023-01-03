package main

import (
	"os"

	"github.com/kondukto-io/simple-fuzzing/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
