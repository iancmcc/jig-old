package main

import (
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "jig"
	app.Usage = "Do something"
	app.Run(os.Args)
}
