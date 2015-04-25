package jig

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
)

type Jigroot interface {
	GetJigroot() string
}

type Jig struct {
	path  string
	repos []*Repository
}

func NewJig(path string) (*Jig, error) {
	return &Jig{path: path, repos: []*Repository{}}, nil
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
func (j *Jig) Reconcile(r Repository) error {
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
