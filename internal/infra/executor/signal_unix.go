//go:build unix

package executor

import (
	"os"
	"syscall"
)

var Signals = []os.Signal{
	syscall.SIGINT,
	syscall.SIGTERM,
}
