package commands

import (
	"sync"

	log "github.com/Sirupsen/logrus"

	"github.com/iancmcc/jig/git"
	"github.com/iancmcc/jig/jig"
	"github.com/jessevdk/go-flags"
)

var _ flags.Commander = &Get{}

func init() {
	parser.AddCommand("get", "Get new repositories into the current Jig environment", "Get new repositories into the current Jig environment", &Get{})
}

type Get struct {
	jigroot_args
	Args struct {
		Repositories []string `positional-arg-name:"REPOSITORY" required:"yes"`
	} `positional-args:"yes" required:"yes"`
}

func (g *Get) getRepo(uri jig.RepositoryURI) error {
	j, err := g.ResolveJig()
	if err != nil {
		return err
	}
	repo := git.GitRepository{jig.BaseRepository{URI: uri}}
	return j.Reconcile(&repo)
}

func (g *Get) Execute(args []string) error {
	g.InitializeLogging()
	if len(g.Args.Repositories) < 1 {
		return &flags.Error{flags.ErrRequired, "Error: At least one repository must be specified"}
	}
	log.Debug("Running get command")
	var wg sync.WaitGroup
	for _, repo := range g.Args.Repositories {
		uri, err := git.ParseGitURI(repo)
		if err != nil {
			return err
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.WithFields(log.Fields{
				"repository": uri,
			}).Info("Getting repository")
			g.getRepo(uri)
		}()
	}
	wg.Wait()
	return nil
}