package executor

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Executor interface {
	ExecCommand(name string, arg ...string) error
	OutputCommand(name string, arg ...string) ([]byte, error)
}

type DefaultExecutor struct{}

func (e *DefaultExecutor) ExecCommand(name string, arg ...string) error {
	fmt.Printf("Executing: %s %s\n", name, strings.Join(arg, " "))
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (e *DefaultExecutor) OutputCommand(name string, arg ...string) ([]byte, error) {
	fmt.Printf("Executing: %s %s\n", name, strings.Join(arg, " "))
	cmd := exec.Command(name, arg...)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("command failed: %w\nstderr:\n%s", err, stderr.String())
	}

	return out, nil
}
