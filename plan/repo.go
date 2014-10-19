package plan

import (
	"fmt"
	"net/url"
	"strings"
)

// Repo represents a version control repository.
type Repo struct {
	FullName string  `json:"name"`
	RefSpec  string  `json:"ref"`
	UserType vcsType `json:"type,omitempty"`
	URI      string  `json:"uri,omitempty"`
	vcsType  string
}

func NewRepo(uri string) (Repo, error) {
	var repo Repo
	parsed, err := url.Parse(uri)
	if err != nil {
		return repo, err
	}
	if strings.Contains(parsed.Path, ":") {
		// It's probably an Git SSH URL
		// TODO: What else could it be? Be defensive
		split := strings.Split(uri, ":")
		nuri := "git+ssh://" + split[0] + "/" + split[1]
		parsed, err = url.Parse(nuri)
		if err != nil {
			return repo, err
		}
	}
	repo = Repo{
		FullName: fmt.Sprintf("%s/%s", parsed.Host, strings.Trim(parsed.Path, "/")),
		URI:      uri,
	}
	return repo, nil
}

func splitFullName(fullname string) ([]string, error) {
	split := strings.SplitN(fullname, "/", 3)
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

// VCSType returns the VCS type of the repository. It will autodetect from
// the repository's registry, or use a user-specified one.
func (r *Repo) VCSType() (vcsType, error) {
	return guessVCSType(r)
}
