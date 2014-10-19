package plan

import "fmt"

type vcsType string

const (
	GIT vcsType = "git"        // git
	SVN vcsType = "subversion" // subversion
)

// ParseRepoKey takes the key of a repo, which may be a URL or a path-style
// name, and determines the missing information.
func ParseRepoKey(raw string) (url, path string, t vcsType, err error) {
	return "", "", GIT, nil
}

func guessVCSType(r *Repo) (vcsType, error) {
	if r.UserType != "" {
		return r.UserType, nil
	}
	reg, err := r.Registry()
	if err != nil {
		return "", err
	}
	if reg == "github.com" {
		return GIT, nil
	}
	return "", fmt.Errorf("Unable to guess repository type from name %s", r.FullName)
}
