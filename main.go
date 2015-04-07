package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/iancmcc/jig/commands"
)

func init() {
	log.SetOutput(os.Stderr)
	log.SetLevel(log.InfoLevel)
}

func main() {
	err := commands.Execute()
	log.WithFields(log.Fields{
		"exitcode": err.ExitCode,
	}).Debug("Exiting")
	os.Exit(err.ExitCode)
}