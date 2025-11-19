//go:build !unix

package os

import (
	"os"
	"syscall"
)

//nolint:gochecknoglobals
var Signals = []os.Signal{
	syscall.SIGINT,
	syscall.SIGTERM,
}
