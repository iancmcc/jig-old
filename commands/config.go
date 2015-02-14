package commands

import (
	"os"
	"os/user"
	"path/filepath"
)

type JigOptions struct {
	Verbose      []bool `short:"v" long:"verbose" description:"Show verbose debug information"`
	JigConfigDir string `short:"c" long:"config" description:"Jig configuration directory" env:"JIG_CONFIG_DIR" default:"$HOME/.jigconfig"`
	IsTTY        bool
}

/*
* ConfigDir() returns the configured Jig configuration directory, or the default
* of ~/.jigconfig.
 */
func (o *JigOptions) ConfigDir() (string, error) {
	if o.JigConfigDir == "$HOME/.jigconfig" {
		usr, err := user.Current()
		if err != nil {
			return "", err
		}
		o.JigConfigDir = filepath.Join(usr.HomeDir, ".jigconfig")
	}
	return o.JigConfigDir, nil
}

/*
* EnsureConfigDir() is the same as ConfigDir(), but ensures the directory
* exists first.
 */
func (o *JigOptions) EnsureConfigDir() (string, error) {
	d, err := o.ConfigDir()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(d, 0755); err != nil {
		return "", err
	}
	return d, nil
}
