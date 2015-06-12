package commands

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/jessevdk/go-flags"
)

var (
	d *Dir            = &Dir{}
	_ flags.Commander = d
)

func init() {
	parser.AddCommand("dir", "Output Jig directory matching search string", "Output Jig directory matching search string", d)
}

type SearchTerm string

type Dir struct {
	jigroot_args
	Args struct {
		Term SearchTerm
	} `positional-args:"true"`
}

func (t *SearchTerm) Complete(match string) []flags.Completion {
	if err := d.Initialize(); err != nil {
		return nil
	}
	repos, err := d.getSimilarRepositories(match)
	if err != nil {
		return nil
	}
	results := []flags.Completion{}
	for _, s := range repos {
		results = append(results, flags.Completion{Item: s})
	}
	return results
}

func (d *Dir) Execute(args []string) error {
	if err := d.Initialize(); err != nil {
		return err
	}
	if d.Args.Term == "" {
		fmt.Println(d.Jigroot)
		return nil
	}
	repos, err := d.getSimilarRepositories(string(d.Args.Term))
	if err != nil {
		return nil
	}
	for _, s := range repos {
		fmt.Println(s)
		return nil
	}
	return nil
}

func (d *Dir) getSimilarRepositories(term string) ([]string, error) {
	j, err := d.ResolveJig()
	if err != nil {
		return nil, err
	}
	repos := j.ListRepositories()
	result := []string{}
	for _, s := range repos {
		trimmed := strings.TrimPrefix(s, d.Jigroot+"/")
		log.Debug("Checking against %s", trimmed)
		w := []rune(trimmed)
		x, _ := FuzzyMatch(false, &w, []rune(term))
		if x > 0 {
			result = append(result, s)
		}
	}
	return result, err
}