package jig

type Repository interface {
}

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
}

type URIScheme string

type URIParser func(uri string) (RepositoryURI, error)