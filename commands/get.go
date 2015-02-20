package commands

import "github.com/jessevdk/go-flags"

var _ flags.Commander = &Get{}

func init() {
	parser.AddCommand("get", "Get new repositories into the current Jig environment", "Get new repositories into the current Jig environment", &Get{})
}

type Get struct {
	Args struct {
		Repositories []string `positional-arg-name:"REPOSITORY" required:"yes"`
	} `positional-args:"yes" required:"yes"`
}

func (g *Get) Execute(args []string) error {
	if len(g.Args.Repositories) < 1 {
		return &flags.Error{flags.ErrRequired, "Error: At least one repository must be specified"}
	}
	return nil
}