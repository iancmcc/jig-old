package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/iancmcc/jig/plan"
	"github.com/iancmcc/jig/vcs"
	"github.com/iancmcc/jig/workbench"
)

func Initialize(ctx *cli.Context) {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Unable to get current working directory")
		os.Exit(1)
	}
	bench := workbench.NewWorkbench(pwd)
	bench.Initialize()
	repo := &plan.Repo{FullName: "github.com/zenoss/platform-build"}
	bench.AddRepository(repo)

	srcrepo, _ := vcs.NewSourceRepository(repo, bench.SrcRoot())
	srcrepo.Create()
}

func main() {
	app := cli.NewApp()
	app.Name = "jig"
	app.Usage = "Do something"
	init := cli.Command{
		Name:   "init",
		Usage:  "Initialize a workbench",
		Action: Initialize,
	}
	app.Commands = []cli.Command{init}
	app.Run(os.Args)
}
