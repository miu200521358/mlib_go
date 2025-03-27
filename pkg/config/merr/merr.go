package merr

import "errors"

var NameNotFoundError = errors.New("name not found")

var TerminateError = errors.New("terminate error")
