//go:build !unix

package os

import (
	"os/exec"
)

type Executor interface {
	Exec(argv0 string, argv []string, envv []string) error
}

type executor struct {
}

func NewExecutor() Executor {
	return &executor{}
}

func (e *executor) Exec(argv0 string, argv []string, envv []string) error {
	cmd := exec.Command(argv0, argv[1:]...)
	cmd.Env = envv

	return cmd.Run()
}
