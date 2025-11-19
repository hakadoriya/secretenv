//go:build !unix

package os

import (
	"os"
	"syscall"
)

var Signals = []os.Signal{
	syscall.SIGINT,
	syscall.SIGTERM,
}
