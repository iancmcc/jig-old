package jig

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/iancmcc/jig/repository"
)

type Jigroot interface {
	GetJigroot() string
}

type Jig struct {
	path  string
	repos []*repository.Repository
}

func NewJig(path string) (*Jig, error) {
	return &Jig{path: path, repos: []*repository.Repository{}}, nil
}

func FindClosestJig(path string) (j *Jig, err error) {
	if path, err = filepath.Abs(path); err != nil {
		return nil, err
	}
	for path != "." && path != "/" {
		if j, err := NewJig(path); err != nil {
			log.Error("Unable to find jig")
			return nil, err
		} else if j.IsRoot() {
			log.WithFields(log.Fields{
				"path": path,
			}).Debug("Found Jigroot")
			return j, nil
		}
		path = filepath.Dir(path)
	}
	return nil, fmt.Errorf("No jig found")
}

func (j *Jig) Path() string {
	return j.path
}

func (j *Jig) jigMetaPath() string {
	return filepath.Join(j.Path(), ".jig")
}

func (j *Jig) Initialize() error {
	if j.IsRoot() {
		return nil
	}
	return os.Mkdir(j.jigMetaPath(), os.ModeDir|0755)
}

func (j *Jig) IsRoot() bool {
	if _, err := os.Stat(j.jigMetaPath()); err != nil {
		return false
	}
	return true
}

/*
* Reconcile() manifests a repository state within a Jig.
 */
func (j *Jig) Reconcile(r repository.Repository) error {
	r.SetRoot(j.Path())
	if err := r.EnsurePath(); err != nil {
		return err
	}
	if err := r.Clone(); err != nil {
		return err
	}
	return nil
}

func (j *Jig) ReconcileAll() error {
	//r.SetRoot(j.Path())
	return nil
}

func (j *Jig) ListRepositories() []repository.Repository {
	// TODO: Replace this with a cache with checksum verification
	repos := []repository.Repository{}
	repoChecker := func(path string, info os.FileInfo, err error) error {
		// Check for a directory named .git
		_, e := os.Stat(filepath.Join(path, ".git"))
		if e == nil {
			// If found, append to repos and return SkipDir
			repo, err := repository.GitRepositoryFromPath(path)
			if err != nil {
				return err
			}
			repos = append(repos, &repo)
			return filepath.SkipDir
		}
		return nil
	}
	filepath.Walk(j.Path(), repoChecker)
	return repos
}

func (j *Jig) repositoryMap() map[string]repository.Repository {
	result := make(map[string]repository.Repository)
	for _, repo := range j.ListRepositories() {
		// Get the path relative to the root
		trimmed := strings.TrimPrefix(repo.GetRoot(), j.path+"/")
		result[trimmed] = repo
	}
	return result
}

func (j *Jig) FindRepositories(term string) []repository.Repository {
	repos := j.repositoryMap()
	result := []repository.Repository{}
	candidates := []string{}
	for k, _ := range repos {
		candidates = append(candidates, k)
	}
	for _, s := range SortedMatches(term, candidates) {
		result = append(result, repos[s])
	}
	return result
}
