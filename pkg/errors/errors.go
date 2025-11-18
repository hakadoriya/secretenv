package errors

import "errors"

var (
	ErrUnknownProvider = errors.New("unknown provider")
	ErrNoArguments     = errors.New("no arguments")
)
