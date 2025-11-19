//go:build unix

package executor

import (
	"fmt"
	"os/exec"

	"golang.org/x/sys/unix"
)

type Executor interface {
	Exec(argv0 string, argv []string, envv []string) error
}

type executor struct{}

func NewExecutor() Executor {
	return &executor{}
}

func (e *executor) Exec(argv0 string, argv []string, envv []string) error {
	execPath, err := exec.LookPath(argv0)
	if err != nil {
		return fmt.Errorf("exec.LookPath: %w", err)
	}

	if err := unix.Exec(execPath, argv, envv); err != nil {
		return fmt.Errorf("unix.Exec: %w", err)
	}

	return nil
}
