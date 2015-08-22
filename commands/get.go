package commands

import (
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/jessevdk/go-flags"
	"github.com/iancmcc/jig/repository"
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

func (g *Get) getRepo(uri repository.RepositoryURI) error {
	repo := repository.GitRepository{repository.BaseRepository{URI: uri}}
	return curJig.Reconcile(&repo)
}

func (g *Get) Execute(args []string) error {
	if err := g.Initialize(); err != nil {
		return err
	}
	if len(g.Args.Repositories) < 1 {
		return &flags.Error{flags.ErrRequired, "Error: At least one repository must be specified"}
	}
	log.Debug("Running get command")
	var wg sync.WaitGroup
	for _, repo := range g.Args.Repositories {
		uri, err := repository.ParseGitURI(repo)
		if err != nil {
			return err
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			t := time.Now()
			g.getRepo(uri)
			log.WithFields(log.Fields{
				"repository": uri,
				"time":       time.Since(t).Seconds(),
			}).Debug("Getting repository")
		}()
	}
	wg.Wait()
	return nil
}