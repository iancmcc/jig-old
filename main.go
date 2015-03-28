package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/iancmcc/jig/commands"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.Info("Welcome to Jig")
	err := commands.Execute()
	log.WithFields(log.Fields{
		"exitcode": err.ExitCode,
	}).Debug("Exiting")
	os.Exit(err.ExitCode)
}