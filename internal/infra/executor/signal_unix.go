//go:build unix

package executor

import (
	"os"
	"syscall"
)

//nolint:gochecknoglobals
var Signals = []os.Signal{
	syscall.SIGINT,
	syscall.SIGTERM,
}
