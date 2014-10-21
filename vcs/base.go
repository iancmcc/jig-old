package vcs

import (
	"path/filepath"

	"github.com/iancmcc/jig/plan"
)

type SourceRepository interface {
	Create() error
}

func NewSourceRepository(repo *plan.Repo, srcroot string, bank *ProgressBarBank) (SourceRepository, error) {
	switch t, _ := repo.VCSType(); t {
	case plan.GIT:
		return NewGitRepository(repo.FullName, filepath.Join(srcroot, repo.FullName), repo.URI, bank), nil
	}
	return nil, nil
}