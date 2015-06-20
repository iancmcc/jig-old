package repository

import (
	"fmt"
	"os"
	"path/filepath"
)

type URIScheme string

const (
	SCHEME_HTTPS URIScheme = "https"
	SCHEME_SSH             = "ssh"
	SCHEME_GIT             = "git"
	SCHEME_EMPTY           = ""
)

type RepositoryURI interface {
	Scheme() URIScheme
	Domain() string
	Owner() string
	Repository() string
	Path() string
	String() string
}

type Repository interface {
	// GetURI() returns the URI describing this repository
	GetURI() *RepositoryURI
	// EnsurePath() creates and returns a filesystem path
	EnsurePath() error
	// Clone() fetches a repository from a remote service
	Clone() error
	// SetRoot() allows a root path to be set
	SetRoot(prefix string)
	// SetRoot() allows a root path to be set
	GetRoot() string
	// GetPath() gets the full path, including the root, returning an error if
	// no root has been set
	GetPath() (string, error)
}

type BaseRepository struct {
	URI  RepositoryURI
	root string
}

func (r *BaseRepository) GetURI() *RepositoryURI {
	return &r.URI
}

func (r *BaseRepository) EnsurePath() error {
	p, err := r.GetPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(p, os.ModeDir|0755); err != nil {
		return err
	}
	return nil
}

func (r *BaseRepository) SetRoot(prefix string) {
	r.root = prefix
}

func (r *BaseRepository) GetRoot() string {
	return r.root
}

func (r *BaseRepository) GetPath() (string, error) {
	if r.root == "" {
		return "", fmt.Errorf("No root has yet been set")
	}
	uri := *r.GetURI()
	return filepath.Join(r.root, uri.Path()), nil
}
