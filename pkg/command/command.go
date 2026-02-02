package command

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Executor interface {
	ExecCommand(name string, arg ...string) error
}

type DefaultExecutor struct{}

func (e *DefaultExecutor) ExecCommand(name string, arg ...string) error {
	fmt.Printf("Executing: %s %s\n", name, strings.Join(arg, " "))
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
