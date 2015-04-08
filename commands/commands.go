package commands

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"

	"github.com/iancmcc/jig/jig"
	"github.com/jessevdk/go-flags"
)

var (
	parser  *flags.Parser = flags.NewNamedParser("jig", flags.Default)
	options JigOptions
)

var curJig *jig.Jig

type jigroot_args struct {
	Jigroot string `short:"j" long:"jigroot" description:"Path to Jigroot" env:"JIGROOT"`
	Init    bool   `long:"init" description:"Initialize a Jigroot"`
	Verbose bool   `short:"v" long:"verbose" description:"Display verbose logging"`
}

func (j *jigroot_args) ResolveJig() (*jig.Jig, error) {
	if j.Jigroot == "" {
		if jj, err := jig.FindClosestJig(j.Jigroot); err == nil {
			log.WithFields(log.Fields{
				"path": jj.Path(),
			}).Info("Found jig root")
			return jj, err
		}
	}
	p, err := filepath.Abs(j.Jigroot)
	if err != nil {
		return nil, err
	}
	jj, err := jig.NewJig(p)
	if !jj.IsRoot() && !j.Init {
		return nil, fmt.Errorf("No Jigroot found. If you want to create one, pass --init")
	}
	jj.Initialize()
	return jj, nil
}

func (j *jigroot_args) Initialize() error {
	var err error
	if j.Verbose {
		log.SetLevel(log.DebugLevel)
	}
	if curJig, err = j.ResolveJig(); err != nil {
		return err
	}
	return nil
}

func Execute() ExecError {
	if _, err := parser.AddGroup("Jig Options", "Jig options", &options); err != nil {
		os.Exit(1)
	}
	_, err := parser.Parse()
	return newExecError(err)
}