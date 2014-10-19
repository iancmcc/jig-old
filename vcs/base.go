package vcs

import (
	"fmt"
	"path/filepath"

	"github.com/iancmcc/jig/plan"
)

type SourceRepository interface {
	Create() error
}

func NewSourceRepository(repo *plan.Repo, srcroot string) (SourceRepository, error) {
	switch t, _ := repo.VCSType(); t {
	case plan.GIT:
		fmt.Printf("Printing repo URI: %s\n", repo.URI)
		return NewGitRepository(filepath.Join(srcroot, repo.FullName), repo.URI), nil
	}
	return nil, nil
}