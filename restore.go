package main

import "fmt"

type RestoreOptions struct {
	Args struct {
		Jigfile   string
		Directory string
	} `positional-args:""`
}

var restoreOptions *RestoreOptions

func init() {
	restoreOptions = &RestoreOptions{}
	parser.AddCommand("restore", "Restore a Jigfile", "Restore a Jigfile to a directory, cloning all necessary repos", restoreOptions)
}

func (r *RestoreOptions) Execute(args []string) error {
	fmt.Printf("%+v", args)
	return nil
}
