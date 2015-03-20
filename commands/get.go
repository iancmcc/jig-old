package commands

import (
	"fmt"
	"sync"

	"github.com/iancmcc/jig/git"
	"github.com/iancmcc/jig/jig"
	"github.com/jessevdk/go-flags"
)

var _ flags.Commander = &Get{}

func init() {
	parser.AddCommand("get", "Get new repositories into the current Jig environment", "Get new repositories into the current Jig environment", &Get{})
}

type Get struct {
	Args struct {
		Repositories []string `positional-arg-name:"REPOSITORY" required:"yes"`
	} `positional-args:"yes" required:"yes"`
}

func (g *Get) getRepo(uri jig.RepositoryURI) {
	fmt.Println(uri.Path())
}

func (g *Get) Execute(args []string) error {
	if len(g.Args.Repositories) < 1 {
		return &flags.Error{flags.ErrRequired, "Error: At least one repository must be specified"}
	}
	var wg sync.WaitGroup
	for _, repo := range g.Args.Repositories {
		uri, err := git.ParseGitURI(repo)
		if err != nil {
			return err
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			g.getRepo(uri)
		}()
	}
	wg.Wait()
	return nil
}