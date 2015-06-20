package repository

import (
	"fmt"
	"regexp"
)

var (
	patterns map[URIScheme]*regexp.Regexp = map[URIScheme]*regexp.Regexp{
		SCHEME_HTTPS: regexp.MustCompile(`https://(?P<domain>.+)/(?P<owner>.+)/(?P<repo>.+)(\.git)?`),
		SCHEME_SSH:   regexp.MustCompile(`git@(?P<domain>.+):(?P<owner>.+)/(?P<repo>.+)(\.git)?`),
		SCHEME_GIT:   regexp.MustCompile(`git://(?P<domain>.+)/(?P<owner>.+)/(?P<repo>.+)(\.git)?`),
		SCHEME_EMPTY: regexp.MustCompile(`((?P<domain>.+)/)?(?P<owner>.+)/(?P<repo>.+)(\.git)?`),
	}
	scheme_order []URIScheme = []URIScheme{SCHEME_SSH, SCHEME_HTTPS, SCHEME_GIT, SCHEME_EMPTY}
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

type URIParser func(uri string) (RepositoryURI, error)

// Type checking
var (
	_ URIParser     = ParseGitURI
	_ RepositoryURI = &GitURI{}
)

type GitURI struct {
	domain string
	owner  string
	repo   string
	scheme URIScheme
}

func ParseGitURI(uri string) (RepositoryURI, error) {
	var (
		domain  string
		owner   string
		repo    string
		scheme  URIScheme
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
			if domain == "" {
				domain = "github.com"
			}
			owner = result["owner"]
			repo = result["repo"]
			break
		}
	}
	if repo == "" {
		// No match
		return nil, fmt.Errorf("Unable to parse git URI: %s", uri)
	}
	if scheme == SCHEME_EMPTY {
		scheme = SCHEME_HTTPS
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

func (u *GitURI) Scheme() URIScheme {
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

func (u *GitURI) Path() string {
	return fmt.Sprintf("%s/%s/%s", u.domain, u.owner, u.repo)
}

func (u *GitURI) String() string {
	switch u.scheme {
	case SCHEME_HTTPS:
		return u.ToHTTPS()
	case SCHEME_SSH:
		return u.ToSSH()
	case SCHEME_GIT:
		return u.ToGit()
	default:
		return u.ToHTTPS()
	}
}