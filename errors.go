package gar

import "fmt"

type ErrWithCmdOutput struct {
	Underlying error
	Output     *output
}

func (e *ErrWithCmdOutput) Error() string {
	return fmt.Sprintf("error while running command: %v; %s", e.Underlying, e.Output.Stderr())
}
