package plan

import (
	"fmt"
	"strings"
)

type Repo struct {
	FullName string `json:"name"`
}

func splitFullName(fullname string) ([]string, error) {
	split := strings.Split(fullname, "/")
	if len(split) != 3 {
		return nil, fmt.Errorf("%s is an invalid repository name", fullname)
	}
	return split, nil
}

// Registry parses FullName and returns the name of the registry (e.g.,
// "github.com" for a repo named "github.com/iancmcc/jig")
func (r *Repo) Registry() (string, error) {
	split, err := splitFullName(r.FullName)
	if err != nil {
		return "", err
	}
	return split[0], nil
}

// Owner parses FullName and returns the name of the owner (e.g.,
// "iancmcc" for a repo named "github.com/iancmcc/jig")
func (r *Repo) Owner() (string, error) {
	split, err := splitFullName(r.FullName)
	if err != nil {
		return "", err
	}
	return split[1], nil
}

// Repository parses FullName and returns the name of the repository (e.g.,
// "jig" for a repo named "github.com/iancmcc/jig")
func (r *Repo) Repository() (string, error) {
	split, err := splitFullName(r.FullName)
	if err != nil {
		return "", err
	}
	return split[2], nil
}
