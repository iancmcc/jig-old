package vcs

import (
	"fmt"
	"os"

	"github.com/libgit2/git2go"
)

// Verify that GitRepository satisfies the SourceRepository interface
var _ SourceRepository = &GitRepository{}

type GitRepository struct {
	site  string
	owner string
	name  string
	path  string
}

func file_exists(filename string) bool {
	_, err := os.Stat(os.ExpandEnv(filename))
	return !os.IsNotExist(err)
}

func getSSHFiles() (string, string, error) {
	dsa_pub := os.ExpandEnv("$HOME/.ssh/id_dsa.pub")
	dsa_key := os.ExpandEnv("$HOME/.ssh/id_dsa")
	rsa_pub := os.ExpandEnv("$HOME/.ssh/id_rsa.pub")
	rsa_key := os.ExpandEnv("$HOME/.ssh/id_rsa")
	if file_exists(dsa_key) && file_exists(dsa_pub) {
		return dsa_pub, dsa_key, nil
	}
	if file_exists(rsa_key) && file_exists(rsa_pub) {
		return rsa_pub, rsa_key, nil
	}
	return "", "", fmt.Errorf("No SSH keys could be found")
}

func credentialsCallback(url string, username_from_url string, allowed_types git.CredType) (int, *git.Cred) {
	pub, key, err := getSSHFiles()
	if err != nil {
		return 0, nil
	}
	i, cred := git.NewCredSshKey("git", pub, key, "")
	return i, &cred
}

func progressCallback(stats git.TransferProgress) int {
	//fmt.Printf("%+v", stats)
	return 0
}

func (r *GitRepository) Create() error {
	// TODO: Allow git to specify its protocol; use SSH only for now
	url := "git@" + r.site + ":" + r.owner + "/" + r.name
	opts := &git.CloneOptions{
		CheckoutOpts: &git.CheckoutOpts{
			Strategy: git.CheckoutSafeCreate,
		},
		RemoteCallbacks: &git.RemoteCallbacks{
			CredentialsCallback: credentialsCallback,
			//TransferProgressCallback: progressCallback,
		},
	}
	_, err := git.Clone(url, r.path, opts)
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	return nil
}