package errors

import "errors"

var (
	ErrEnvVarNotSet    = errors.New("environment variable is not set")
	ErrUnknownProvider = errors.New("unknown provider")
	ErrNoArguments     = errors.New("no arguments")
)
