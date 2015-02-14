package commands

type JigOptions struct {
	Verbose []bool `short:"v" long:"verbose" description:"Show verbose debug information"`

	IsTTY bool
}
