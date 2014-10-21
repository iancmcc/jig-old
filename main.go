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

	r, _ := plan.NewRepo("git@github.com:zenoss/platform-build")
	j, _ := plan.NewRepo("git@github.com:zenoss/manifest")
	k, _ := plan.NewRepo("git@github.com:iancmcc/dotfiles")
	m, _ := plan.NewRepo("git@github.com:control-center/serviced")
	bench.AddRepository(&r)
	bench.AddRepository(&j)
	bench.AddRepository(&k)
	bench.AddRepository(&m)

	var wg sync.WaitGroup
	bank := vcs.NewProgressBarBank()
	for name, r := range bench.Plan().Repos {
		if name != "" {
			wg.Add(1)
			go func(name string, r *plan.Repo) {
				defer wg.Done()
				srcrepo, _ := vcs.NewSourceRepository(r, bench.SrcRoot(), bank)
				if err := srcrepo.Create(); err != nil {
					fmt.Printf("%+v\n", err)
				}
			}(name, r)
		}
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
