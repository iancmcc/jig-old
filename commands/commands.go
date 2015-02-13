package commands

import "fmt"

func Execute() ExecError {
	fmt.Println("run")
	return newExecError(nil)
}