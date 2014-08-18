package vcs

import (
	"path/filepath"

	"github.com/iancmcc/jig/plan"
)

type SourceRepository interface {
	Create() error
}

func NewSourceRepository(repo *plan.Repo, srcroot string) (SourceRepository, error) {
	switch t, _ := repo.VCSType(); t {
	case plan.GIT:
		site, err := repo.Registry()
		if err != nil {
			return nil, err
		}
		owner, err := repo.Owner()
		if err != nil {
			return nil, err
		}
		name, err := repo.Repository()
		if err != nil {
			return nil, err
		}
		return &GitRepository{
			site:  site,
			owner: owner,
			name:  name,
			path:  filepath.Join(srcroot, repo.FullName),
		}, nil

	}
	return nil, nil
}