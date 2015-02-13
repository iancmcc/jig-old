package main

import (
	"os"

	"github.com/iancmcc/jig/commands"
)

func main() {
	err := commands.Execute()
	os.Exit(err.ExitCode)
}