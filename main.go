package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/iancmcc/jig/plan"
	"github.com/iancmcc/jig/workbench"
)

func getWorkbench(ctx *cli.Context) (*workbench.Workbench, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	fmt.Println("Finding nearest bench")
	if bench := workbench.FindNearestBench(pwd); bench != nil {
		return bench, nil
	}
	// We're not in an existing bench, so let's get a plan and create a bench
	// from it
	fmt.Println("No existing bench")
	plan, err := getPlan(ctx)
	if err != nil {
		return nil, err
	}
	return workbench.NewWorkbench(pwd, plan), nil
}

// plan() discovers the current bench plan in the following order:
//   1. Plan specified by first argument
//   3. Plan from Jigfile in cwd
//   4. Empty plan
// First one wins.
func getPlan(ctx *cli.Context) (*plan.Plan, error) {
	// Plan specified by first argument
	fmt.Println("Getting plan")
	filename := ctx.Args().First()
	if filename == "" {
		// Jigfile in current dir
		filename = "Jigfile"
	}
	fmt.Printf("Trying plan filename %s\n", filename)
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		f, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		fmt.Printf("New plan from json")
		return plan.NewPlanFromJSON(f)
	}
	// Empty plan
	fmt.Printf("Empty plan")
	return plan.NewPlan(), nil
}

func Add(ctx *cli.Context) {
}

func Initialize(ctx *cli.Context) {
	bench, err := getWorkbench(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	bench.Realize()
}

func main() {
	app := cli.NewApp()
	app.Name = "jig"
	app.Usage = "Do something"
	init := cli.Command{
		Name:   "apply",
		Usage:  "Apply changes to an existing workbench, or create a new one",
		Action: Initialize,
	}
	add := cli.Command{
		Name:   "add",
		Usage:  "Add a repository",
		Action: Add,
	}
	app.Commands = []cli.Command{init, add}
	app.Run(os.Args)
}
