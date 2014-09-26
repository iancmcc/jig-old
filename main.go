package main

import (
	"fmt"
	"os"
	"sync"

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

	bench.AddRepository(&plan.Repo{FullName: "github.com/control-center/serviced"})
	bench.AddRepository(&plan.Repo{FullName: "github.com/zenoss/platform-build"})
	bench.AddRepository(&plan.Repo{FullName: "github.com/iancmcc/dotfiles"})

	var wg sync.WaitGroup
	for name, r := range bench.Plan().Repos {
		wg.Add(1)
		go func(name string, r *plan.Repo) {
			defer wg.Done()
			fmt.Println("Starting clone of", name)
			srcrepo, _ := vcs.NewSourceRepository(r, bench.SrcRoot())
			srcrepo.Create()
			fmt.Println("Done with", name)
		}(name, r)
	}
	wg.Wait()
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
