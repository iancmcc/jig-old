package vcs

import "github.com/libgit2/git2go"

// Verify that GitRepository satisfies the SourceRepository interface
var _ SourceRepository = &GitRepository{}

type GitRepository struct {
	site  string
	owner string
	name  string
	path  string
}

func (r *GitRepository) Create() error {
	// TODO: Allow git to specify its protocol; use SSH only for now
	url := "https://" + r.site + "/" + r.owner + "/" + r.name + ".git"
	opts := &git.CloneOptions{
		CheckoutOpts: &git.CheckoutOpts{
			Strategy: git.CheckoutSafeCreate,
		},
	}
	git.Clone(url, r.path, opts)
	return nil
}