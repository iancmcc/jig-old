package git

import (
	"fmt"
	"regexp"

	"github.com/iancmcc/jig/jig"
)

var (
	patterns map[jig.URIScheme]*regexp.Regexp = map[jig.URIScheme]*regexp.Regexp{
		jig.SCHEME_HTTPS: regexp.MustCompile(`https://(?P<domain>.+)/(?P<owner>.+)/(?P<repo>.+)(\.git)?`),
		jig.SCHEME_SSH:   regexp.MustCompile(`git@(?P<domain>.+):(?P<owner>.+)/(?P<repo>.+)(\.git)?`),
		jig.SCHEME_GIT:   regexp.MustCompile(`git://(?P<domain>.+)/(?P<owner>.+)/(?P<repo>.+)(\.git)?`),
		jig.SCHEME_EMPTY: regexp.MustCompile(`(?P<domain>.+)/(?P<owner>.+)/(?P<repo>.+)(\.git)?`),
	}
	scheme_order []jig.URIScheme = []jig.URIScheme{jig.SCHEME_SSH, jig.SCHEME_HTTPS, jig.SCHEME_GIT, jig.SCHEME_EMPTY}
)

// Type checking
var (
	_ jig.URIParser     = ParseGitURI
	_ jig.RepositoryURI = &GitURI{}
)

type GitURI struct {
	domain string
	owner  string
	repo   string
	scheme jig.URIScheme
}

func ParseGitURI(uri string) (jig.RepositoryURI, error) {
	var (
		domain  string
		owner   string
		repo    string
		scheme  jig.URIScheme
		pattern *regexp.Regexp
	)
	for _, scheme = range scheme_order {
		pattern = patterns[scheme]
		if pattern.MatchString(uri) {
			match := pattern.FindStringSubmatch(uri)
			result := make(map[string]string)
			for i, name := range pattern.SubexpNames() {
				result[name] = match[i]
			}
			domain = result["domain"]
			owner = result["owner"]
			repo = result["repo"]
			break
		}
	}
	if domain == "" {
		// No match
		return nil, fmt.Errorf("Unable to parse git URI: %s", uri)
	}
	if scheme == jig.SCHEME_EMPTY {
		scheme = jig.SCHEME_HTTPS
	}
	return &GitURI{
		domain: domain,
		scheme: scheme,
		owner:  owner,
		repo:   repo,
	}, nil
}

func (u *GitURI) Domain() string {
	return u.domain
}

func (u *GitURI) Owner() string {
	return u.owner
}

func (u *GitURI) Repository() string {
	return u.repo
}

func (u *GitURI) Scheme() jig.URIScheme {
	return u.scheme
}

func (u *GitURI) ToSSH() string {
	return fmt.Sprintf("git@%s:%s/%s", u.domain, u.owner, u.repo)
}

func (u *GitURI) ToHTTPS() string {
	return fmt.Sprintf("https://%s/%s/%s", u.domain, u.owner, u.repo)
}

func (u *GitURI) ToGit() string {
	return fmt.Sprintf("git://%s/%s/%s", u.domain, u.owner, u.repo)
}