package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

var completion string = `
%s() {
    cd $(jig dir $@)
}

_jig_completion() {
    # All arguments except the first one
    args=("${COMP_WORDS[@]:1:$COMP_CWORD}")

    # Only split on newlines
    local IFS=$'\n'

    # Call completion (note that the first element of COMP_WORDS is
    # the executable itself)
    COMPREPLY=($(GO_FLAGS_COMPLETION=1 ${COMP_WORDS[0]} "${args[@]}"))
    return 0
}

complete -F _jig_completion jig
`

func init() {
	parser.AddCommand("bootstrap", "Bootstrap a shell environment", "Bootstrap a shell environment", &Bootstrap{})
}

type Bootstrap struct {
	Alias string `short:"a" long:"alias" description:"Specify the alias used for changing to source directories" default:"cdj"`
}

func (b *Bootstrap) writeToFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	jigDir := filepath.Join(usr.HomeDir, ".jig")
	// Ensure the directory is created
	os.MkdirAll(jigDir, 0755)

	completed := fmt.Sprintf(completion, b.Alias)

	// Write to a file
	bootfile := filepath.Join(jigDir, "bootstrap~")
	if err := ioutil.WriteFile(bootfile, []byte(completed), 0644); err != nil {
		return "", err
	}
	return bootfile, nil
}

func (b *Bootstrap) Execute(args []string) error {
	path, err := b.writeToFile()
	if err != nil {
		return err
	}
	fmt.Println(path)
	return nil
}