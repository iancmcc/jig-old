package vcs

import (
	"fmt"
	"os"

	"github.com/libgit2/git2go"
)

// Verify that GitRepository satisfies the SourceRepository interface
var _ SourceRepository = &GitRepository{}

type GitRepository struct {
	path string
	url  string
}

func NewGitRepository(path, url string) *GitRepository {
	return &GitRepository{path, url}
}

func file_exists(filename string) bool {
	_, err := os.Stat(os.ExpandEnv(filename))
	return !os.IsNotExist(err)
}

func getSSHFiles() (string, string, error) {
	dsa_pub := os.ExpandEnv("$HOME/.ssh/id_dsa.pub")
	dsa_key := os.ExpandEnv("$HOME/.ssh/id_dsa")
	if file_exists(dsa_key) && file_exists(dsa_pub) {
		fmt.Printf("FILES YO")
		return dsa_pub, dsa_key, nil
	}
	rsa_pub := os.ExpandEnv("$HOME/.ssh/id_rsa.pub")
	rsa_key := os.ExpandEnv("$HOME/.ssh/id_rsa")
	if file_exists(rsa_key) && file_exists(rsa_pub) {
		fmt.Printf("RSA FILES YO")
		return rsa_pub, rsa_key, nil
	}
	return "", "", fmt.Errorf("No SSH keys could be found")
}

func credentialsCallback(url string, username_from_url string, allowed_types git.CredType) (int, *git.Cred) {
	fmt.Printf("url: %s, username_from_url: %s; allowed: %+v", url, username_from_url, allowed_types)
	pub, key, err := getSSHFiles()
	if err != nil {
		return 0, nil
	}
	i, cred := git.NewCredSshKey("git", pub, key, "")
	return i, &cred
}

func certCheckCallback(cert *git.Certificate, valid bool, hostname string) int {
	// Don't bother checking any certs, just go with it if valid
	if valid {
		return 1
	}
	return 0
}

func progressCallback(stats git.TransferProgress) int {
	fmt.Printf("%+v", stats)
	return 0
}

func (r *GitRepository) Create() error {
	// TODO: Allow git to specify its protocol; use SSH only for now
	opts := &git.CloneOptions{
		CheckoutOpts: &git.CheckoutOpts{
			Strategy: git.CheckoutSafeCreate,
		},
		RemoteCallbacks: &git.RemoteCallbacks{
			CertificateCheckCallback: certCheckCallback,
			CredentialsCallback:      credentialsCallback,
			TransferProgressCallback: progressCallback,
		},
	}
	_, err := git.Clone(r.url, r.path, opts)
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	return nil
}