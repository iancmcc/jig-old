package commands

import (
	"os"

	"github.com/jessevdk/go-flags"
)

var (
	parser  *flags.Parser = flags.NewNamedParser("jig", flags.Default)
	options JigOptions
)

func Execute() ExecError {
	if _, err := parser.AddGroup("Jig Options", "Jig options", &options); err != nil {
		os.Exit(1)
	}
	_, err := parser.Parse()
	return newExecError(err)
}